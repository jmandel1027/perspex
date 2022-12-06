package router

import (
	"net/http"

	"github.com/dimiro1/health"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"

	"github.com/jmandel1027/perspex/services/gateway/pkg/config"
	"github.com/jmandel1027/perspex/services/gateway/pkg/handler/graphql"
)

// Route -- used for mounting all of our routes
func Route(cfg *config.GatewayConfig) http.Handler {
	router := http.NewServeMux()

	gql := graphql.GraphQL(cfg)

	router.Handle("/api/graphql", gql)
	router.Handle("/api/metrics", promhttp.Handler())
	router.Handle("/api/status", health.NewHandler())

	return cors.AllowAll().Handler(router)
}
