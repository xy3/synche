package config_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xy3/synche/src/config"
	c "github.com/xy3/synche/src/server"
	"testing"
)

func TestWrite(t *testing.T) {
	err := config.Write("test-config.yaml", c.Config)
	assert.NoError(t, err)
}
