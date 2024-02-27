package logger

import (
	"context"
	"testing"
)

func TestNewLoggerForSmoke(t *testing.T) {
	logger := NewLogger(context.Background(), "test logger")
	logger.Errorf("a test")
}
