package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func DSN(DBUsername string, DBPassword string, DBProtocol string, DBAddress string) (dsn string) {
	return fmt.Sprintf("%s:%s@%s(%s)/", DBUsername, DBPassword, DBProtocol, DBAddress)
}

func CreateDatabase(driverName string, DBUsername string, DBPassword string, DBProtocol string, DBAddress string, name string) (err error) {
	/* Data source name configuration has the following parameters:
	  [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN] */
	dataSourceName := DSN(DBUsername, DBPassword, DBProtocol, DBAddress)

	// Create a DB without connecting to an existing DB, validates DSN
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Printf("Error when opening database: %s", err)
		return err
	}
	log.Printf("Database %s opened successfully", dataSourceName)

	/* Create the DB now that connection has been established
	  Timeout in case of connectivity or runtime errors */
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Returns a pool of underlying DB connections
	result, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		log.Printf("Error when creating database: %s", err)
		return err
	}
	log.Printf("Database %s created successfully", name)

	no, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error when fetching rows: %s", err)
		return err
	}
	log.Printf("Rows affected when creating database: %d", no)

	/* Need to close the existing connection to the DB and open a new
	connection with the correct DB name that we just created */
	if err := db.Close(); err != nil {
		log.Print("Could not close connection to database")
	}

	// Create and connect to database
	db, err = DBConnection(driverName, dataSourceName, name)
	if err != nil {
		log.Printf("Error when getting db connection: %s", err)
		return
	}
	defer db.Close()
	log.Printf("Successfully connected to database")

	// Add chunk table to database by default
	err = CreateChunkTable(db)
	if err != nil {
		log.Printf("Creating chunk table failed with error %s", err)
		return
	}
	log.Printf("Successfully created chunk table")

	// Add connection request table to database by default
	err = CreateConnectionRequestTable(db)
	if err != nil {
		log.Printf("Creating connection_request table failed with error %s", err)
		return
	}
	log.Printf("Successfully created connection_request table")

	if err := db.Close(); err != nil {
		log.Print("Could not close connection to database")
	}
	return nil
}

func DBConnection(driverName string, dataSourceName string, name string) (db *sql.DB, err error){
	db, err = sql.Open(driverName, dataSourceName + name)
	if err != nil {
		log.Printf("Error when opening database: %s", err)
		return nil, err
	}

	// Recommended values for not starving a DB, may need to be reviewed
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 2)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Verify the connection
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Error pinging database: %s", err)
		return nil, err
	}
	log.Printf("Connected to database %s successfully", name)

	return db, nil
}

func CreateChunkTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS chunk(file_name varchar(255), file_size int(64), chunk_hash varchar(255) primary key, chunk_number int(64), upload_request_id varchar(255), file_chunk_directory varchar(255))`

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	result, err := db.ExecContext(ctx, query)

	if err != nil {
		log.Printf("Error when creating chunk table: %s", err)
		return err
	}
	rows, err := result.RowsAffected()

	if err != nil {
		log.Printf("Error when getting rows affected: %s", err)
		return err
	}

	log.Printf("Rows affected when creating chunk table: %d", rows)
	return nil
}

func InsertChunk(db *sql.DB, fileName string, fileSize int64, chunkHash string, chunkNumber int64, uploadRequestId string, directoryId string) error {
	query := "INSERT INTO chunk(file_name, file_size, chunk_hash, chunk_number, upload_request_id, file_chunk_directory) VALUES (?, ?, ?, ?, ?, ?)" +
		"ON DUPLICATE KEY UPDATE file_name=VALUES(file_name), file_size=VALUES(file_size), chunk_number=VALUES(chunk_number), upload_request_id=VALUES(upload_request_id), file_chunk_directory=VALUES(file_chunk_directory)"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error when preparing SQL statement:\n-->%s", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, fileName, fileSize, chunkHash, chunkNumber, uploadRequestId, directoryId)
	if err != nil {
		log.Print("Error when inserting row into chunk table")
		log.Print(err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error when finding rows affected: %s", err)
		return err
	}

	log.Printf("Chunk information inserted in to database for %d chunk(s)", rows)
	return nil
}


func CreateConnectionRequestTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS connection_request(upload_request_id varchar(255) primary key, file_chunk_directory varchar(255), file_name varchar(255), file_size int(64), number_of_chunks int(64))`

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	result, err := db.ExecContext(ctx, query)

	if err != nil {
		log.Printf("Error when creating connection_request table: %s", err)
		return err
	}
	rows, err := result.RowsAffected()

	if err != nil {
		log.Printf("Error when getting rows affected: %s", err)
		return err
	}

	log.Printf("Rows affected when creating table: %d", rows)
	return nil
}

func InsertConnectionRequest(db *sql.DB, uploadRequestId string, fileChunkDir string, fileName string, fileSize int64, numberOfChunks int64) error {
	query := "INSERT INTO connection_request(upload_request_id, file_chunk_directory, file_name, file_size, number_of_chunks) VALUES (?, ?, ?, ?, ?)" +
		"ON DUPLICATE KEY UPDATE file_chunk_directory=VALUES(file_chunk_directory), file_name=VALUES(file_name), file_size=VALUES(file_size), number_of_chunks=VALUES(number_of_chunks)"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error when preparing SQL statement: %s", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, uploadRequestId, fileChunkDir, fileName, fileSize, numberOfChunks)
	if err != nil {
		log.Printf("Error when inserting row into upload_request table: %s", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error when finding rows affected: %s", err)
		return err
	}

	log.Printf("%d request ID added: ", rows)
	return nil
}
