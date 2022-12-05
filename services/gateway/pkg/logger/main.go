package logger

import (
	"context"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// Declare variables to store log messages as new Events
var (
	invalidArgMessage      = Event{1, "Invalid arg: %s"}
	invalidArgValueMessage = Event{2, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{3, "Missing arg: %s"}
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
	logger.Info("logger construction succeeded")

	return otelzap.New(logger)
}

// ReplaceGlobals replaces zap's global loggers with the provided logger.
func ReplaceGlobals(logger *otelzap.Logger) func() {
	undo := otelzap.ReplaceGlobals(logger)
	return undo
}

// InvalidArg is a standard error message
func InvalidArg(ctx context.Context, argumentName string) {
	otelzap.L().Ctx(ctx).Sugar().Errorf(invalidArgMessage.message, argumentName)
}

// InvalidArgValue is a standard error message
func InvalidArgValue(ctx context.Context, argumentName string, argumentValue string) {
	otelzap.L().Ctx(ctx).Sugar().Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

// MissingArg is a standard error message
func MissingArg(ctx context.Context, argumentName string) {
	otelzap.L().Ctx(ctx).Sugar().Errorf(missingArgMessage.message, argumentName)
}

// Info Log
func Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	otelzap.L().Ctx(ctx).Info(msg, fields...)
}

// Warn Log
func Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	otelzap.Ctx(ctx).Sugar().Warnf(msg, fields)
}

// Debug Log
func Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	otelzap.Ctx(ctx).Sugar().Debugf(msg, fields)
}

// Panic Log
func Panic(ctx context.Context, msg string, fields ...zapcore.Field) {
	otelzap.Ctx(ctx).Sugar().Panicf(msg, fields)
}

// Error Log
func Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	otelzap.Ctx(ctx).Sugar().Errorf(msg, fields)
}

// Fatal Log
func Fatal(ctx context.Context, msg string, fields ...zapcore.Field) {
	otelzap.Ctx(ctx).Sugar().Fatalf(msg, fields)
}
