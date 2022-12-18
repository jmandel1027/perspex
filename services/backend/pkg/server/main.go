package server

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cockroachdb/cmux"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	users "github.com/jmandel1027/perspex/schemas/proto/goproto/pkg/users/v1"
	config "github.com/jmandel1027/perspex/services/backend/pkg/config"
	"github.com/jmandel1027/perspex/services/backend/pkg/database/postgres"
	"github.com/jmandel1027/perspex/services/backend/pkg/logger"
	grpc_postgres "github.com/jmandel1027/perspex/services/backend/pkg/middleware/postgres"
	"github.com/jmandel1027/perspex/services/backend/pkg/router"
	"github.com/jmandel1027/perspex/services/backend/pkg/tracing"
	userService "github.com/jmandel1027/perspex/services/backend/pkg/user/service"
)

func Serve() {

	cfg, err := config.New()
	if err != nil {
		otelzap.L().Warn("Config Error: %s", zap.Error(err))
	}

	z := logger.New(cfg)
	defer z.Sync()

	undo := logger.ReplaceGlobals(z)
	defer undo()

	l, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		otelzap.L().Fatal("Listen Error: %s", zap.Error(err))
	}

	m := cmux.New(l)
	go func() { m.Serve() }()

	grpc := m.Match(cmux.HTTP2())

	go GRPC(&cfg, grpc)
	// go HTTP(&cfg)

	select {}

}

// GRPC -- Handler
func GRPC(cfg *config.BackendConfig, l net.Listener) {
	dbs, err := postgres.Open(cfg)
	if err != nil {
		otelzap.L().Warn("Postgres Connection Error: %s", zap.Error(err))
	}

	trace, err := tracing.NewTracerProvider(context.TODO())
	if err != nil {
		otelzap.L().Ctx(context.TODO()).Error("Tracer Provider Error: %s", zap.Error(err))
	}

	defer func() {
		if err := trace.Shutdown(context.Background()); err != nil {
			otelzap.L().Ctx(context.TODO()).Error("Error shutting down tracer provider")
		}
	}()

	middleware := grpc_middleware.ChainUnaryServer(
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(otelzap.L().Ctx(context.TODO()).ZapLogger()),
		grpc_recovery.UnaryServerInterceptor(),
		grpc_postgres.UnaryServerInterceptor(dbs.Writer, &sql.TxOptions{ReadOnly: false}),
		grpc_postgres.UnaryServerInterceptor(dbs.Reader, &sql.TxOptions{ReadOnly: true}),
	)

	unary := grpc.UnaryInterceptor(middleware)
	sever := grpc.NewServer(unary)

	userService := userService.NewUserService()
	users.RegisterUserServiceServer(sever, userService)

	reflection.Register(sever)

	if err := sever.Serve(l); err != nil {
		otelzap.Ctx(context.TODO()).Panic("gRPC Serve Error: %s", zap.Error(err))
	}

	otelzap.L().Info("gRPC Server Stopped")

	defer func() {
		// Here is where we'd safely close out any connections
		// eg: redis, grpc, etc..
	}()

}

// HTTP server
func HTTP(cfg *config.BackendConfig) {
	ctx := context.Background()

	otelzap.L().Ctx(ctx).Info("Scaffolded global logger")

	rtr := router.Route(cfg)

	srv := &http.Server{
		Addr:         cfg.Host + ":" + cfg.HttpPort,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 10,
		Handler:      rtr,
	}

	// Create a goroutine that listens for interupts
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Here is a non blocking go routine that runs forever
		// it's our listener that exposes the entire app
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			otelzap.L().Ctx(ctx).Fatal("Error:", zap.String("err", err.Error()))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 60000*time.Second)

	// when the shutoff signal is received, the lock will release
	// and anything below done will run.
	<-done

	otelzap.L().Ctx(ctx).Info("Backend Server Stopped")

	defer func() {
		// Here is where we'd safely close out any connections
		// eg: redis, etc. Except for Postgres, we need to allow that package to manage
		// It's own lifecycle.

		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		otelzap.L().Ctx(ctx).Fatal("Server Shutdown Failed:", zap.String("err", err.Error()))
	}

	otelzap.L().Ctx(ctx).Info("Server Exited Properly")

	defer os.Exit(0)
}
