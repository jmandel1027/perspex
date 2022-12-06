package logger

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a new logger
func New() *otelzap.Logger {
	config := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			TimeKey:       "time",
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			StacktraceKey: "stack",
			CallerKey:     "caller",
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
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
