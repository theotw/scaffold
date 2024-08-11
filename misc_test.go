package scaffold

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MiscTestSuite struct {
	suite.Suite
}

func TestStuff(t *testing.T) {
	suite.Run(t, new(MiscTestSuite))
}

func (suite *MiscTestSuite) TestGetEnv() {
	suite.Assert().Equal(GetEnv("TEST", "default"), "default")
	os.Setenv("TEST", "test")
	suite.Assert().Equal(GetEnv("TEST", "default"), "test")
}
