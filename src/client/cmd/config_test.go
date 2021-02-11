package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"io/ioutil"
	"testing"
)

var (
	testConfigFile = "testdata/.config.yaml"
)

func TestConfigCommand(t *testing.T) {
	// Use this global log hook to get the last log entry
	hook := test.NewGlobal()
	// Avoid printing any log output during the tests
	log.SetOutput(ioutil.Discard)

	configCmd := NewConfigCommand()

	err := config.InitConfig(testConfigFile)
	if err != nil {
		t.Error(err)
	}

	testCases := []struct {
		Name     string
		Args     []string
		Expected string
	}{
		{"check output is the same as viper config file", []string{}, "Config file path: " + testConfigFile},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v ", tc.Name), func(t *testing.T) {
			configCmd.SetArgs(tc.Args)

			if err := configCmd.Execute(); err != nil {
				t.Error(err)
			}

			assert.Equal(t, 1, len(hook.Entries))
			assert.Equal(t, tc.Expected, hook.LastEntry().Message)
		})
	}
}
