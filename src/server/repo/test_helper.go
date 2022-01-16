package repo

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xy3/synche/src/files"
	schema2 "github.com/xy3/synche/src/schema"
	"github.com/xy3/synche/src/server"
	"gorm.io/gorm"
	"testing"
)

// NewUserForTest sets up a new test database, creates a user and a home directory for them
func NewUserForTest(t *testing.T) (
	user *schema2.User,
	homeDir *schema2.Directory,
	db *gorm.DB,
	down func(*testing.T),
	err error,
) {
	db, down = server.NewTxForTest(t)
	user, err = newTestUser(db)
	require.Nil(t, err)

	homeDir, err = SetupUserHomeDir(user, db)
	assert.Nil(t, err)

	return user, homeDir, db, down, err
}

func newTestUser(db *gorm.DB) (*schema2.User, error) {
	testUser := server.TestUser
	user := &schema2.User{
		Email:     testUser.Email,
		EmailHash: files.MD5HashString(testUser.Email),
		Password:  testUser.Password,
		TokenHash: files.MD5HashString(testUser.PlaintextPassword),
	}

	tx := db.Where(user).FirstOrCreate(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
