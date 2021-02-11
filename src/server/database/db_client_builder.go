package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type DBClientBuilder struct {
	// this place can be used to initialize the auth client
}

func NewDBClientBuilder() DBClientBuilder {
	return DBClientBuilder{}
}

func (b DBClientBuilder) BuildSqlClient(driver string, username string, password string, protocol string, address string, name string) *sql.DB {
	// sensitive info can be stored in "secrets.json" of GKE
	dataSourceName := DSN(username, password, protocol, address)

	db, err := sql.Open(driver, dataSourceName + name)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err.Error())
	}
	return db
}