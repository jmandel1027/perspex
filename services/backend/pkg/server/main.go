package server

import (
	"log"
	"net"

	"github.com/cockroachdb/cmux"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	users "github.com/jmandel1027/perspex/schemas/proto/goproto/pkg/users/v1"
	"github.com/jmandel1027/perspex/services/backend/pkg/config"
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

	svr := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(),
		)),
	)

	userService := userService.NewUserService()
	users.RegisterUserServiceServer(svr, userService)

	reflection.Register(svr)

	if err := svr.Serve(l); err != nil {
		log.Panicf("[Server.GRPC] grpc serve error: %s", err)
	}

	defer func() {
		// Here is where we'd safely close out any connections
		// eg: redis, grpc, etc..
	}()
}
