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

	"github.com/jmandel1027/perspex/services/gateway/pkg/config"
	"github.com/jmandel1027/perspex/services/gateway/pkg/router"
)

// HTTP server
func HTTP(cfg *config.GatewayConfig) {
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
