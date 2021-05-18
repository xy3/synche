package ftp

import (
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
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

	correct, err = auth.CheckPasswd(database.TestUser.Email, database.TestUser.PlaintextPassword)
	assert.Nil(t, err)
	assert.True(t, correct)
}
