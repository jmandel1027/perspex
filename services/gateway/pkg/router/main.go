package router

import (
	"net/http"

	"github.com/jmandel1027/perspex/services/gateway/pkg/config"

	"github.com/rs/cors"

	"github.com/dimiro1/health"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Route -- used for mounting all of our routes
func Route(cfg *config.GatewayConfig) http.Handler {
	rtr := http.NewServeMux()

	rtr.Handle("/api/metrics", promhttp.Handler())
	rtr.Handle("/api/status", health.NewHandler())

	return cors.AllowAll().Handler(rtr)
}
