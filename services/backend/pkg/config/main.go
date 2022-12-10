package config

import (
	"fmt"
	"time"

	"github.com/jmandel1027/perspex/services/backend/pkg/utils"
)

// LogConfig defines the configuration for the logger
type LogConfig struct {
	Verbose bool
}

// PostgresConfig configures the PostgreSQL connection.
type PostgresConfig struct {
	Name         string
	User         string
	Password     string
	Host         string
	Port         string
	Schema       string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifespan  int
	Debug        bool
}

// RedisConfig defines a redis connection configuration.
type RedisConfig struct {
	Host               string
	Password           string
	Port               string
	DB                 int
	InsecureSkipVerify bool
	TLS                bool
}

// BackendConfig defines the configuration for the server
type BackendConfig struct {
	Host     string
	HttpPort string
	GrpcPort string
	Log      LogConfig
	WriterPG PostgresConfig
	ReaderPG PostgresConfig
	Redis    RedisConfig
}

// New return all constants using in Project
func New() (BackendConfig, error) {
	log := LogConfig{
		Verbose: utils.MustGetBool("BACKEND_LOG_MODE", "false"),
	}

	writerPG := PostgresConfig{
		Name:         utils.MustGet("WRITER_POSTGRES_DB", "perspex"),
		User:         utils.MustGet("WRITER_POSTGRES_USER", "perspex"),
		Password:     utils.MustGet("WRITER_POSTGRES_PASSWORD", "pass"),
		Host:         utils.MustGet("WRITER_POSTGRES_HOST", "postgresql.perspex.svc.cluster.local"),
		Port:         utils.MustGet("WRITER_POSTGRES_PORT", "5432"),
		Schema:       utils.MustGet("WRITER_POSTGRES_SCHEMA", "public"),
		MaxOpenConns: utils.MustGetInt("WRITER_POSTGRES_MAXOPENCONNECTIONS", "100"),
		MaxIdleConns: utils.MustGetInt("WRITER_POSTGRES_MAXOPENCONNECTIONS", "50"),
		MaxLifespan:  utils.MustGetInt("WRITER_POSTGRES_CONNECTION_LIFESPAN", "128"),
		Debug:        utils.MustGetBool("WRITER_POSTGRES_DEBUG", "false"),
	}

	readerPG := PostgresConfig{
		Name:         utils.MustGet("READER_POSTGRES_DB", "perspex"),
		User:         utils.MustGet("READER_POSTGRES_USER", "perspex"),
		Password:     utils.MustGet("READER_POSTGRES_PASSWORD", "pass"),
		Host:         utils.MustGet("READER_POSTGRES_HOST", "postgresql.perspex.svc.cluster.local"),
		Port:         utils.MustGet("READER_POSTGRES_PORT", "5432"),
		Schema:       utils.MustGet("READER_POSTGRES_SCHEMA", "public"),
		MaxOpenConns: utils.MustGetInt("READER_POSTGRES_MAXOPENCONNECTIONS", "100"),
		MaxIdleConns: utils.MustGetInt("READER_POSTGRES_MAXOPENCONNECTIONS", "50"),
		MaxLifespan:  utils.MustGetInt("READER_POSTGRES_CONNECTION_LIFESPAN", "128"),
		Debug:        utils.MustGetBool("READER_POSTGRES_DEBUG", "false"),
	}

	redis := RedisConfig{
		Host:               utils.MustGet("REDIS_HOST", "redis-writer.perspex"),
		Password:           utils.MustGet("REDIS_PASSWORD", "pass"),
		Port:               utils.MustGet("REDIS_PORT", "6369"),
		DB:                 utils.MustGetInt("REDIS_DB", "0"),
		InsecureSkipVerify: utils.MustGetBool("REDIS_SKIP_VERIFY", "false"),
		TLS:                utils.MustGetBool("REDIS_TLS", "false"),
	}

	return BackendConfig{
		Host:     utils.MustGet("BACKEND_HOST", "0.0.0.0"),
		HttpPort: utils.MustGet("BACKEND_HTTP_PORT", "8080"),
		GrpcPort: utils.MustGet("BACKEND_GRPC_PORT", "8888"),
		Log:      log,
		WriterPG: writerPG,
		ReaderPG: readerPG,
		Redis:    redis,
	}, nil
}

// GetDebug returns true if debug information should be outputted to the DebugWriter handler.
func (d PostgresConfig) GetDebug() bool {
	return d.Debug
}

// GetDataSourceName constructs a postgres DSN.
func (d PostgresConfig) GetDataSourceName() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.Host,
		d.Port,
		d.User,
		d.Password,
		d.Name,
	)
}

// GetMaxOpenConns gets the maximum number of open connections to the database.
func (d PostgresConfig) GetMaxOpenConns() int {
	return d.MaxOpenConns
}

// GetMaxIdleConns gets the maximum number of connections in the idle connection pool.
func (d PostgresConfig) GetMaxIdleConns() int {
	return d.MaxIdleConns
}

// GetConnMaxLifetime gets the maximum amount of time a connection may be reused.
func (d PostgresConfig) GetConnMaxLifetime() time.Duration {
	return time.Duration(d.MaxLifespan) * time.Second
}
