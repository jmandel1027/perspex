package server

import (
	"context"
	"log"
	"net"

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
	"github.com/jmandel1027/perspex/services/backend/pkg/tracing"
	userService "github.com/jmandel1027/perspex/services/backend/pkg/user/service"
)

func Serve() {
	cfg, err := config.New()
	if err != nil {
		log.Panic(err)
	}

	l, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		log.Fatal(err)
	}

	m := cmux.New(l)
	go func() { m.Serve() }()

	grpc := m.Match(cmux.HTTP2())

	go GRPC(&cfg, grpc)

	select {}

}

// GRPC -- Handler
func GRPC(cfg *config.BackendConfig, l net.Listener) {
	trace, err := tracing.NewTracerProvider()
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
		grpc_zap.UnaryServerInterceptor(otelzap.L().Logger),
		grpc_recovery.UnaryServerInterceptor(),
	)

	unary := grpc.UnaryInterceptor(middleware)
	sever := grpc.NewServer(unary)

	userService := userService.NewUserService()
	users.RegisterUserServiceServer(sever, userService)

	reflection.Register(sever)

	if err := sever.Serve(l); err != nil {
		otelzap.L().Ctx(context.TODO()).Panic("gRPC Serve Error: %s", zap.Error(err))
	}

	otelzap.L().Ctx(context.TODO()).Info("gRPC Server Stopped")

	defer func() {
		// Here is where we'd safely close out any connections
		// eg: redis, grpc, etc..
	}()

}
