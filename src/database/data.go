package database

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	DB *gorm.DB
)

func migrateAll(db *gorm.DB) error {
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

	var err error
	DB, err = gorm.Open(mysql.Open(c.Config.Database.DSN()), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		log.WithError(err).Error("Failed to open gorm DB connection")
		return DB, err
	}

	return configureConnection(DB)
}

// RequireNewConnection just calls NewConnection, but it will panic when something goes wrong
func RequireNewConnection() *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)
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
	if err := sqlDB.PingContext(ctx); err != nil {
		return db, err
	}
	return db, nil
}

func InitSyncheData() error {
	db, err := NewConnection()
	if err != nil {
		return err
	}

	if err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", c.Config.Database.Name)).Error; err != nil {
		return err
	}

	if err = migrateAll(db); err != nil {
		log.WithError(err).Error("Failed to migrate database")
		return err
	}

	return nil
}
