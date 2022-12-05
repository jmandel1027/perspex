package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/jmandel1027/perspex/schemas/graphql/pkg/resolvers"
	"github.com/jmandel1027/perspex/schemas/graphql/pkg/source"
	"github.com/jmandel1027/perspex/services/gateway/pkg/config"
)

// GraphQL --
func GraphQL(cfg *config.GatewayConfig) http.HandlerFunc {

	resolvers
	c := source.Config{
		Resolvers: &resolvers.Resolver{},
	}

	h := handler.New(source.NewExecutableSchema(c))

	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{})
	h.Use(extension.Introspection{})

	return h.ServeHTTP
}
