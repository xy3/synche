package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"log"
)

type DBClientBuilder struct {
	// this place can be used to initialize the auth client
}

func NewDBClientBuilder() DBClientBuilder {
	return DBClientBuilder{}
}

func (b DBClientBuilder) BuildSqlClient(dbConfig config.DatabaseConfig) *sql.DB {
	// sensitive info can be stored in "secrets.json" of GKE
	dataSourceName := NewDSN(dbConfig)

	db, err := sql.Open(dbConfig.Driver, dataSourceName + dbConfig.Name)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err.Error())
	}
	return db
}