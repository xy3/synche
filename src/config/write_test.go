package config_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"testing"
)

func TestWrite(t *testing.T) {
	err := config.Write("test-config.yaml", c.Config)
	assert.NoError(t, err)
}
