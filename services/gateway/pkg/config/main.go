package config

import (
	"github.com/jmandel1027/perspex/services/gateway/pkg/logger"
	"github.com/jmandel1027/perspex/services/gateway/pkg/utils"
)

// RedisConfig defines a redis connection configuration.
type RedisConfig struct {
	Host               string
	Password           string
	Port               string
	DB                 int
	InsecureSkipVerify bool
	TLS                bool
}

type UserService struct {
	GrpcPort string
}

// GatewayConfig defines the configuration for the server
type GatewayConfig struct {
	Host     string
	HttpPort string
	GrpcPort string
	Redis    RedisConfig
}

// New return all constants using in Project
func New() (GatewayConfig, error) {

	z := logger.New()
	defer z.Sync()

	undo := logger.ReplaceGlobals(z)
	defer undo()

	return GatewayConfig{
		Host:     utils.MustGet("BACKEND_HOST", "0.0.0.0"),
		HttpPort: utils.MustGet("BACKEND_HTTP_PORT", "8080"),
		GrpcPort: utils.MustGet("BACKEND_GRPC_PORT", "8888"),
		Redis: RedisConfig{
			Host:               utils.MustGet("REDIS_HOST", "redis-writer.perspex"),
			Password:           utils.MustGet("REDIS_PASSWORD", "pass"),
			Port:               utils.MustGet("REDIS_PORT", "6369"),
			DB:                 0,
			InsecureSkipVerify: utils.MustGetBool("REDIS_SKIP_VERIFY", "false"),
			TLS:                utils.MustGetBool("REDIS_TLS", "false"),
		},
	}, nil
}
