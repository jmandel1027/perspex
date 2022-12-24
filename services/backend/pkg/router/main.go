package router

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/bufbuild/connect-go"
	grpcReflect "github.com/bufbuild/connect-grpcreflect-go"
	otelconnect "github.com/bufbuild/connect-opentelemetry-go"
	"github.com/dimiro1/health"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	usersConnect "github.com/jmandel1027/perspex/schemas/proto/goproto/pkg/users/v1/usersconnect"
	"github.com/jmandel1027/perspex/services/backend/pkg/config"
	"github.com/jmandel1027/perspex/services/backend/pkg/database/postgres"
	transaction "github.com/jmandel1027/perspex/services/backend/pkg/middleware/postgres"
	userService "github.com/jmandel1027/perspex/services/backend/pkg/user/service"
)

// Route -- used for mounting all of our routes
func Route(cfg *config.BackendConfig, dbs *postgres.DB) http.Handler {
	mux := http.NewServeMux()
	api := http.NewServeMux()

	users := userService.NewUserService()

	reflector := grpcReflect.NewStaticReflector(
		"user.v1.UserService",
	)

	otelzap.L().Info("Scaffolding reader")
	reader := func(ctx context.Context, req *transaction.Request) (*sql.DB, *sql.TxOptions, error) {
		return transaction.Wrap(ctx, dbs.Reader, postgres.ReadOnlyTxOpts)
	}

	otelzap.L().Info("Scaffolding writer")
	writer := func(ctx context.Context, req *transaction.Request) (*sql.DB, *sql.TxOptions, error) {
		return transaction.Wrap(ctx, dbs.Writer, postgres.StdTxOpts)
	}

	otelzap.L().Info("Scaffolding opts")

	opts := connect.WithInterceptors(
		otelconnect.NewInterceptor(),
		transaction.New(reader),
		transaction.New(writer),
	)

	api.Handle(grpcReflect.NewHandlerV1(reflector))
	api.Handle(grpcReflect.NewHandlerV1Alpha(reflector))
	api.Handle(usersConnect.NewUserServiceHandler(users, opts))

	mux.Handle("/", api)
	mux.Handle("/api/metrics", promhttp.Handler())
	mux.Handle("/api/status", health.NewHandler())

	return cors.AllowAll().Handler(mux)
}
