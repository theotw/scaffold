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
	tmpfile := os.TempDir() + "/test.log"
	defer os.Remove(tmpfile)
	f, _ := os.OpenFile(tmpfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	orgStdErr := os.Stderr
	os.Stderr = f
	ctx := NewLoggerContext(context.Background(), "my-tracer")
	os.Setenv(loggerTypeEnvVar, loggerTypeLogrus)
	logger := NewLogger(ctx, "test logrus logger")
	suite.NotNil(logger)
	suite.dumpLogLevels(logger, "logrus")

	os.Setenv(loggerTypeEnvVar, loggerTypeZap)
	logger = NewLogger(ctx, "test zap logger")
	suite.NotNil(logger)
	suite.dumpLogLevels(logger, "zap")
	os.Stderr = orgStdErr
	f.Close()
	fileBits, err := os.ReadFile(tmpfile)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(fileBits)
	suite.Require().NoError(err)
	stderrString := string(fileBits)
	suite.Contains(stderrString, "This is a debug message logrus")
	suite.Contains(stderrString, "This is an info message logrus")
	suite.Contains(stderrString, "This is an error message logrus")
	suite.Contains(stderrString, "This is a debug message zap")
	suite.Contains(stderrString, "This is an info message zap")
	suite.Contains(stderrString, "This is an error message zap")
	// check for the tracer
	suite.Contains(stderrString, "my-tracer")
	// check for the logger names
	suite.Contains(stderrString, "test logrus logger")
	suite.Contains(stderrString, "test zap logger")

}
func (suite *LoggerTestSuite) dumpLogLevels(log Logger, name string) {
	log.Debugf("This is a debug message %s", name)
	log.Infof("This is an info message %s", name)
	log.Errorf("This is an error message %s", name)
}
