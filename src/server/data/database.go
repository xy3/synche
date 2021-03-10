package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
)

type Database struct {
	mysql *sql.DB
}

func NewDatabaseClient(dbConfig config.DatabaseConfig) *Database {
	db, err := sql.Open(dbConfig.Driver, NewDSN(dbConfig) + dbConfig.Name)
	if err != nil {
		log.WithError(err).Fatal("Error connecting to the database")
	}
	return &Database{mysql: db}
}
