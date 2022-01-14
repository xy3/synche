package server

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
)

var (
	TestUser = struct {
		Email             string
		Password          string
		PlaintextPassword string
		Name              string
		Picture           string
	}{
		Email:             "synche@synche.test",
		Password:          "$2y$10$OGO3hi/0uiJI6TvnHX8H2OBAHhmS/NmN66ITHPuG9IBRUCE0P4Smu",
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
	assert.NoError(t, err)

	err = MigrateAll(db)
	assert.NoError(t, err)

	testDB = db
	return testDB
}

func NewTxForTest(t *testing.T) (db *gorm.DB, down func(*testing.T)) {
	defer func() { assert.Nil(t, recover()) }()

	db = newDBForTest(t)
	TestDB = db

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
