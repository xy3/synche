package ftp

import (
	"github.com/stretchr/testify/assert"
	"github.com/xy3/synche/src/server"
	"github.com/xy3/synche/src/server/repo"
	"testing"
)

func TestAuth_CheckPasswd(t *testing.T) {
	_, _, db, down, err := repo.NewUserForTest(t)
	assert.Nil(t, err)
	defer down(t)

	auth := &Auth{testDB: db}

	correct, err := auth.CheckPasswd("", "")
	assert.NotNil(t, err)
	assert.False(t, correct)

	correct, err = auth.CheckPasswd("wrongUsername", "wrongPassword")
	assert.NotNil(t, err)
	assert.False(t, correct)

	correct, err = auth.CheckPasswd(server.TestUser.Email, server.TestUser.PlaintextPassword)
	assert.Nil(t, err)
	assert.True(t, correct)
}
