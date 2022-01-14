package server

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xy3/synche/src/server/schema"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var (
	DB     *gorm.DB
	TestDB *gorm.DB
)

func MigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&schema.Chunk{},
		&schema.File{},
		&schema.FileChunk{},
		&schema.Directory{},
		&schema.User{},
	)
}

func NewConnection() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}

	if TestDB != nil {
		DB = TestDB
		return configureConnection(DB)
	}

	var err error
	DB, err = gorm.Open(mysql.Open(Config.Database.DSN()), &gorm.Config{
		PrepareStmt: true,
		Logger: logger.New(
			log.New(),
			logger.Config{
				IgnoreRecordNotFoundError: true,
			},
		),
	})

	if err != nil {
		log.WithError(err).Error("Failed to open gorm DB connection")
		return DB, err
	}

	return configureConnection(DB)
}

// RequireNewConnection just calls NewConnection, but it will panic when something goes wrong
func RequireNewConnection() (db *gorm.DB) {
	var err error
	if db, err = NewConnection(); err != nil {
		panic(err)
	}
	return db
}

func configureConnection(db *gorm.DB) (*gorm.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}

	// Recommended values for not starving a Database, may need to be reviewed
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	// Verify the connection
	if err = sqlDB.PingContext(ctx); err != nil {
		return db, err
	}
	return db, nil
}

func InitSyncheData() (*gorm.DB, error) {
	db, err := NewConnection()
	if err != nil {
		return nil, err
	}

	if TestDB == nil {
		if err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", Config.Database.Name)).Error; err != nil {
			return nil, err
		}
	}

	if err = MigrateAll(db); err != nil {
		log.WithError(err).Error("Failed to migrate database")
		return nil, err
	}

	return db, nil
}
