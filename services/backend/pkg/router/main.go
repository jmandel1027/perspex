package router

import (
	"net/http"

	"github.com/dimiro1/health"
	"github.com/jmandel1027/perspex/services/backend/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

// Route -- used for mounting all of our routes
func Route(cfg *config.BackendConfig) http.Handler {
	router := http.NewServeMux()

	router.Handle("/api/metrics", promhttp.Handler())
	router.Handle("/api/status", health.NewHandler())

	return cors.AllowAll().Handler(router)
}
