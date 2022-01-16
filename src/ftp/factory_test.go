package ftp

import (
	"github.com/stretchr/testify/assert"
	c "github.com/xy3/synche/src/server"
	"testing"
)

func TestDriverFactory_NewDriver(t *testing.T) {
	assert.Nil(t, c.InitConfig(""))
	_, err := (&Factory{}).NewDriver()
	assert.Nil(t, err)
}
