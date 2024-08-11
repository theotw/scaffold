package logger

import (
	"context"
	"os"
)

type loggerContextKey string

const (
	// loggerTypeEnvVar is the environment variable that will be used to determine the logger type
	loggerTypeEnvVar = "LOGGER_TYPE"
	// loggerTypeZap is the zap logger type
	loggerTypeZap = "zap"
	// loggerTypeLogrus is the logrus logger type (default)
	loggerTypeLogrus = "logrus"
	//TracerField used to follow calls in a call chain
	TracerField = loggerContextKey("tracer")
	//LoggerName used to identify the logger
	LoggerName = "name"
)

// NewLoggerContext creates a new context with the tracer value, this tracer is added to every log message
// used in this context
func NewLoggerContext(ctx context.Context, tracerValue string) context.Context {
	return context.WithValue(ctx, TracerField, tracerValue)
}

// Logger the logging interface the app will use.  Keeping it simple.
// Other loggers pretty much implement these and if needed a wrapper class can be used.
type Logger interface {
	// Infof logs a message at level Info
	Infof(template string, args ...any)
	// Debugf logs a message at level Warn
	Debugf(template string, args ...any)
	// Errorf logs a message at level Error
	Errorf(template string, args ...any)
	// Fatalf logs a message at level Fatal level and then panics
	Fatalf(template string, args ...any)
	// Sync flushes the logger
	Sync() error
}

func NewLogger(ctx context.Context, loggerName string) Logger {
	loggerType := os.Getenv(loggerTypeEnvVar)
	if loggerType == "" {
		loggerType = loggerTypeLogrus
	}
	switch loggerType {
	case loggerTypeZap:
		return NewZapLogger(ctx, loggerName)
	case loggerTypeLogrus:
		return NewLogrusLogger(ctx, loggerName)
	default:
		return NewLogrusLogger(ctx, loggerName)
	}
}
