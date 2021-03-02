package data

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"time"
)

func createTables(db *sql.DB) (err error) {
	// Add chunk table to data by default
	err = CreateChunkTable(db)
	if err != nil {
		return err
	}
	log.Printf("Successfully created chunk table")

	// Add connection request table to data by default
	err = CreateConnectionRequestTable(db)
	if err != nil {
		return err
	}
	log.Info("Successfully created connection_request table")

	return nil
}

func DBConnection(driverName string, dataSourceName string, name string) (db *sql.DB, err error){
	db, err = sql.Open(driverName, dataSourceName + name)
	if err != nil {
		return db, err
	}

	// Recommended values for not starving a Database, may need to be reviewed
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 2)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Verify the connection
	err = db.PingContext(ctx)
	if err != nil {
		return db, err
	}
	log.Infof("Connected to data %s successfully", name)

	return db, nil
}

func NewDSN(dbConfig config.DatabaseConfig) (dsn string) {
	return fmt.Sprintf("%s:%s@%s(%s)/", dbConfig.Username, dbConfig.Password, dbConfig.Protocol, dbConfig.Address)
}

func CreateDatabase(dbConfig config.DatabaseConfig) (err error) {
	/* Data source name configuration has the following parameters:
	  [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN] */
	dataSourceName := NewDSN(dbConfig)

	// Create a Database without connecting to an existing Database, validates NewDSN
	db, err := sql.Open(dbConfig.Driver, dataSourceName)
	if err != nil {
		return err
	}
	log.Infof("Database %s opened successfully", dataSourceName)

	/* Create the Database now that connection has been established
	  Timeout in case of connectivity or runtime errors */
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Returns a pool of underlying Database connections
	result, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS " + dbConfig.Name)
	if err != nil {
		return err
	}
	log.Infof("Database %s created successfully", dbConfig.Name)

	no, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof("Rows affected when creating data: %d", no)

	/* Need to close the existing connection to the Database and open a new
	connection with the correct Database name that we just created */
	if err := db.Close(); err != nil {
		return err
	}

	// Create and connect to data
	db, err = DBConnection(dbConfig.Driver, dataSourceName, dbConfig.Name)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Info("Successfully connected to data")

	if err = createTables(db); err != nil {
	return err
	}
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

func CreateChunkTable(db *sql.DB) (err error) {
	query := `CREATE TABLE IF NOT EXISTS chunk(file_name varchar(255), file_size int(64), chunk_hash varchar(255) primary key, chunk_number int(64), upload_request_id varchar(255), file_chunk_directory varchar(255))`

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Infof("Rows affected when creating chunk table: %d", rows)
	return nil
}

func CreateConnectionRequestTable(db *sql.DB) (err error) {
	query := `CREATE TABLE IF NOT EXISTS connection_request(upload_request_id varchar(255) primary key, file_chunk_directory varchar(255), file_name varchar(255), file_size int(64), number_of_chunks int(64))`

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Infof("Rows affected when creating table: %d", rows)
	return nil
}

