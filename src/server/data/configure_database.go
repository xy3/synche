package data

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"time"
)

func (d *SyncheData) Configure() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return configureDB(sqlDB)
}

func configureDB(sqlDB *sql.DB) error {
	// Recommended values for not starving a Database, may need to be reviewed
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	// Verify the connection
	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func NewDSN(dbConfig config.DatabaseConfig) (dsn string) {
	return fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Protocol,
		dbConfig.Address,
		dbConfig.Name,
	)
}
