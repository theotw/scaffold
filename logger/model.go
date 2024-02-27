package logger

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (

	//TracerField used to follow calls in a call chain
	TracerField = "tracer"
	LoggerName  = "name"
)

// Logger the logging interface the app will use.  Keeping it simple.
// Other loggers pretty much implement these and if needed a wrapper class can be used.
type Logger interface {
	Infof(template string, args ...any)
	Debugf(template string, args ...any)
	Errorf(template string, args ...any)
	Fatalf(template string, args ...any)
	Sync() error
}

func NewLogger(ctx context.Context, loggerName string) Logger {
	tracerValOb := ctx.Value(TracerField)
	var tracerVal string
	if tracerValOb == nil {
		tracerVal = uuid.NewString()
	} else {
		tracerVal = tracerValOb.(string)
	}
	tracer := zap.String(TracerField, tracerVal)
	loggerNameField := zap.String(LoggerName, loggerName)
	fields := zap.Fields(tracer, loggerNameField)
	logger, err := zap.NewProduction(fields)
	if err != nil {
		panic("Go not create the logger")
	}
	sugar := logger.Sugar()
	return sugar
}
