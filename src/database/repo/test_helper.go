package repo

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"gorm.io/gorm"
	"testing"
)

// NewUserForTest sets up a new test database, creates a user and a home directory for them
func NewUserForTest(t *testing.T) (
	user *schema.User,
	homeDir *schema.Directory,
	db *gorm.DB,
	down func(*testing.T),
	err error,
) {
	db, down = database.NewTxForTest(t)
	user, err = newTestUser(db)
	require.Nil(t, err)

	homeDir, err = SetupUserHomeDir(user, db)
	assert.Nil(t, err)

	return user, homeDir, db, down, err
}

func newTestUser(db *gorm.DB) (*schema.User, error) {
	testUser := database.TestUser
	user := &schema.User{
		Email:     testUser.Email,
		EmailHash: hash.MD5HashString(testUser.Email),
		Password:  testUser.Password,
		TokenHash: hash.MD5HashString(testUser.PlaintextPassword),
	}

	tx := db.Where(user).FirstOrCreate(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}