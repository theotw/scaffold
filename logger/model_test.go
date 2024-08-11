package logger

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
}

func TestLogger(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (suite *LoggerTestSuite) SetupTest() {
}
func (suite *LoggerTestSuite) TeardownTest() {

}
func (suite *LoggerTestSuite) TestNewLoggerForSmoke() {
	ctx := NewLoggerContext(context.Background(), "my-tracer")
	os.Setenv(loggerTypeEnvVar, loggerTypeLogrus)
	logger := NewLogger(ctx, "test logrus logger")
	suite.NotNil(logger)
	suite.dumpLogLevels(logger)

	os.Setenv(loggerTypeEnvVar, loggerTypeZap)
	logger = NewLogger(ctx, "test zap logger")
	suite.NotNil(logger)
	suite.dumpLogLevels(logger)
}
func (suite *LoggerTestSuite) dumpLogLevels(log Logger) {
	log.Debugf("This is a debug message %d", 1)
	log.Infof("This is an info message %d", 2)
	log.Errorf("This is an error message %d", 4)
}
