package logger

import (
	"github.com/jmandel1027/perspex/services/backend/pkg/config"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a new logger
func New(cfg config.BackendConfig) *otelzap.Logger {

	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			MessageKey:    "message",
			LevelKey:      "level",
			StacktraceKey: "stack",
			TimeKey:       "time",
		},
	}

	// If verbose logging is enabled, add caller and path fields
	if cfg.Log.Verbose {
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.FunctionKey = "path"
	}

	logger, _ := config.Build()
	logger.Info("Logger construction succeeded")

	return otelzap.New(logger)
}

// ReplaceGlobals replaces zap's global loggers with the provided logger.
func ReplaceGlobals(logger *otelzap.Logger) func() {
	undo := otelzap.ReplaceGlobals(logger)
	return undo
}
