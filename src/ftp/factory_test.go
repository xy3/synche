package ftp

import (
	"github.com/stretchr/testify/assert"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"testing"
)

func TestDriverFactory_NewDriver(t *testing.T) {
	assert.Nil(t, c.InitConfig(""))
	_, err := (&Factory{}).NewDriver()
	assert.Nil(t, err)
}
