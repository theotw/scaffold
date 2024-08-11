package logger

import (
	"context"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Logger = (*zap.SugaredLogger)(nil)

func NewZapLogger(ctx context.Context, loggerName string) Logger {
	ll := os.Getenv("LOG_LEVEL")
	if ll == "" {
		ll = "DEBUG"
	}
	level, err := zapcore.ParseLevel(ll)
	if err != nil {
		level = zapcore.DebugLevel
	}
	tracerValOb := ctx.Value(TracerField)
	var tracerVal string
	if tracerValOb == nil {
		tracerVal = uuid.NewString()
	} else {
		tracerVal = tracerValOb.(string)
	}
	tracer := zap.String(string(TracerField), tracerVal)
	loggerNameField := zap.String(LoggerName, loggerName)
	levelOption := zap.IncreaseLevel(level)
	fields := zap.Fields(tracer, loggerNameField)
	atomicLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)

	// Create the logger core with DebugLevel or lower to ensure all levels are allowed
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stderr),
		atomicLevel, // This allows dynamic level adjustment
	)
	logger := zap.New(core, levelOption, fields)

	sugar := logger.Sugar()
	return sugar
}
