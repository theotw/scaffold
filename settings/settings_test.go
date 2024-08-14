package settings

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type testSettings struct {
	Item1 string `json:"item_1"`
	Item2 string `json:"item_2"`
	Item3 string `json:"item_3"`
	// we skip item 4 on purpose, its an item we load but not in here to test unmarshal
	Item5 string `json:"item_5"`
	Item6 string `json:"item_6"`
}

type testSettingsSource struct {
	suite.Suite
}

func TestSettings(t *testing.T) {
	suite.Run(t, &testSettingsSource{})
}
func (s *testSettingsSource) TestEnvSettings() {
	pwd, _ := os.Getwd()
	envFileName := path.Join(pwd, "testdata", "test1.env")
	_, err := os.Stat(envFileName)
	s.Require().NoError(err)
	settings := NewEnvSettingsSource[testSettings]("OTW_", envFileName)
	err = settings.InitSettings()
	s.Require().NoError(err)
	s.Require().NoError(err)
	s.Require().NotNil(settings.GetSettings())
	s.Assert().Equal("item1", settings.GetItem("item_1"))
	s.Assert().Equal("item2", settings.GetItem("item_2"))
	s.Assert().Equal("item3", settings.GetItem("item_3"))
	s.Assert().Equal("item4", settings.GetItem("item_4"))

	s.Assert().Equal("item1", settings.GetSettings().Item1)
	s.Assert().Equal("item2", settings.GetSettings().Item2)
	s.Assert().Equal("item3", settings.GetSettings().Item3)
}

func (s *testSettingsSource) TestEnvAndSecretsSettings() {
	pwd, _ := os.Getwd()
	envFileName := path.Join(pwd, "testdata", "test2.env")
	_, err := os.Stat(envFileName)
	s.Require().NoError(err)
	settings := NewEnvSettingsSource[testSettings]("OTW_", envFileName)
	err = settings.InitSettings()
	s.Require().NoError(err)
	s.Require().NoError(err)
	s.Require().NotNil(settings.GetSettings())
	s.Assert().Equal("item1", settings.GetItem("item_1"))
	s.Assert().Equal("item2", settings.GetItem("item_2"))
	s.Assert().Equal("item3", settings.GetItem("item_3"))
	s.Assert().Equal("item4", settings.GetItem("item_4"))
	// from aws secret
	s.Assert().Equal("item5", settings.GetItem("item_5"))
	//overriden from aws secret
	s.Assert().Equal("item6-override", settings.GetItem("item_6"))

	s.Assert().Equal("item1", settings.GetSettings().Item1)
	s.Assert().Equal("item2", settings.GetSettings().Item2)
	s.Assert().Equal("item3", settings.GetSettings().Item3)
	s.Assert().Equal("item5", settings.GetSettings().Item5)
	s.Assert().Equal("item6-override", settings.GetSettings().Item6)
}
