package logger

import (
	"context"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var _ Logger = (*LogrusLogger)(nil)
var log *logrus.Logger
var logLock = &sync.Mutex{}

var _ Logger = (*LogrusLogger)(nil)

// LogrusLogger is an implementation of the Logger interface using logrus
type LogrusLogger struct {
	logrus.Entry
}

func (l *LogrusLogger) Sync() error {
	// no-op in this logger
	return nil
}

// NewLogrusLogger creates a new LogrusLogger instance
func NewLogrusLogger(ctx context.Context, name string) Logger {
	logLock.Lock()
	if log == nil {
		log = mkLogger()
	}
	defer logLock.Unlock()

	tracerValOb := ctx.Value(TracerField)
	var tracerVal string
	if tracerValOb == nil {
		tracerVal = uuid.NewString()
	} else {
		tracerVal = tracerValOb.(string)
	}

	withContext := log.WithContext(ctx).WithField(string(TracerField), tracerVal).WithField(LoggerName, name)
	l := &LogrusLogger{*withContext}
	return l
}

// mkLogger creates a new logrus logger instance
func mkLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "DEBUG"
	}
	ll, err := logrus.ParseLevel(level)
	if err != nil {
		ll = logrus.DebugLevel
	}
	logger.SetLevel(ll)
	return logger
}
