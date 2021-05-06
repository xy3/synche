package repo

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
)

var (
	NewUserForTest = newUserForTest

	TestUser = struct {
		Email             string
		PlaintextPassword string
		Name              string
		Picture           string
	}{
		Email:             "synche@synche.test",
		PlaintextPassword: "testPassword123",
		Name:              "testName",
		Picture:           "testPicture",
	}
)

var testDB *gorm.DB

func newDBForTest(t *testing.T) *gorm.DB {
	if testDB != nil {
		return testDB
	}

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	assert.Nil(t, err)

	testDB = db
	return testDB
}

func newTxForTest(t *testing.T) (*gorm.DB, func(*testing.T)) {
	assert.Nil(t, c.InitConfig(""))

	defer func() { assert.Nil(t, recover()) }()

	db := newDBForTest(t)

	database.TestDB = db
	db, err := database.InitSyncheData()
	assert.Nil(t, err)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true,
		},
	)
	db.Logger = newLogger
	tx := db.Begin(&sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	return tx, func(t *testing.T) {
		defer func() { assert.Nil(t, recover()) }()
		tx.Rollback()
	}
}

func newUserForTest(t *testing.T) (
	user *schema.User,
	homeDir *schema.Directory,
	db *gorm.DB,
	down func(*testing.T),
	err error,
) {
	db, down = newTxForTest(t)
	user, err = NewUser(TestUser.Email, TestUser.PlaintextPassword, &TestUser.Name, &TestUser.Picture, db)
	require.Nil(t, err)

	homeDir, err = SetupUserHomeDir(user, db)
	assert.Nil(t, err)

	return user, homeDir, db, down, err
}
