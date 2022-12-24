package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	config "github.com/jmandel1027/perspex/services/backend/pkg/config"
	"github.com/jmandel1027/perspex/services/backend/pkg/database/postgres"
	"github.com/jmandel1027/perspex/services/backend/pkg/logger"
	"github.com/jmandel1027/perspex/services/backend/pkg/router"
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

	dbs, err := postgres.Open(&cfg)
	if err != nil {
		otelzap.L().Warn("Postgres Connection Error: %s", zap.Error(err))
	}

	go HTTP(&cfg, dbs)

	select {}
}

// HTTP server
func HTTP(cfg *config.BackendConfig, dbs *postgres.DB) {
	ctx := context.Background()

	otelzap.L().Ctx(ctx).Info("Scaffolded global logger")

	rtr := router.Route(cfg, dbs)

	srv := &http.Server{
		Addr:           cfg.Host + ":" + cfg.HttpPort,
		WriteTimeout:   time.Minute * 5,
		ReadTimeout:    time.Minute * 5,
		IdleTimeout:    time.Minute * 5,
		MaxHeaderBytes: 8 * 1024, // 8KiB
		Handler:        h2c.NewHandler(rtr, &http2.Server{}),
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
		dbs.Writer.Close()
		dbs.Reader.Close()

		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		otelzap.L().Ctx(ctx).Fatal("Server Shutdown Failed:", zap.String("err", err.Error()))
	}

	otelzap.L().Ctx(ctx).Info("Server Exited Properly")

	defer os.Exit(0)
}
