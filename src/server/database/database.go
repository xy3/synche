package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
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
		return err
	}
	log.Infof("Database %s opened successfully", dataSourceName)

	/* Create the DB now that connection has been established
	  Timeout in case of connectivity or runtime errors */
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	// Returns a pool of underlying DB connections
	result, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		return err
	}
	log.Infof("Database %s created successfully", name)

	no, err := result.RowsAffected()
	if err != nil {
		return err
	}
	log.Infof("Rows affected when creating database: %d", no)

	/* Need to close the existing connection to the DB and open a new
	connection with the correct DB name that we just created */
	if err := db.Close(); err != nil {
		return err
	}

	// Create and connect to database
	db, err = DBConnection(driverName, dataSourceName, name)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Info("Successfully connected to database")

	// Add chunk table to database by default
	err = CreateChunkTable(db)
	if err != nil {
		return err
	}
	log.Printf("Successfully created chunk table")

	// Add connection request table to database by default
	err = CreateConnectionRequestTable(db)
	if err != nil {
		return err
	}
	log.Info("Successfully created connection_request table")

	if err := db.Close(); err != nil {
		return err
	}
	return nil
}

func DBConnection(driverName string, dataSourceName string, name string) (db *sql.DB, err error){
	db, err = sql.Open(driverName, dataSourceName + name)
	if err != nil {
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
		return nil, err
	}
	log.Infof("Connected to database %s successfully", name)

	return db, nil
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

func InsertChunk(db *sql.DB, fileName string, fileSize int64, chunkHash string, chunkNumber int64, uploadRequestId string, directoryId string) (err error) {
	query := "INSERT INTO chunk(file_name, file_size, chunk_hash, chunk_number, upload_request_id, file_chunk_directory) VALUES (?, ?, ?, ?, ?, ?)" +
		"ON DUPLICATE KEY UPDATE file_name=VALUES(file_name), file_size=VALUES(file_size), chunk_number=VALUES(chunk_number), upload_request_id=VALUES(upload_request_id), file_chunk_directory=VALUES(file_chunk_directory)"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, fileName, fileSize, chunkHash, chunkNumber, uploadRequestId, directoryId)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// ON DUPLICATE KEY UPDATE will result in 2 rows being updated if there is a duplicate key
	if rows == 0 {
		log.Warnf("Chunk information already exists in database for chunk: %s", chunkHash)
	} else if rows == 1 {
		log.Infof("Chunk information inserted into database for chunk: %s", chunkHash)
	} else {
		log.Infof("Chunk information updated in database for chunk: %s", chunkHash)
	}
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

func InsertConnectionRequest(db *sql.DB, uploadRequestId string, fileChunkDir string, fileName string, fileSize int64, numberOfChunks int64) (err error) {
	query := "INSERT INTO connection_request(upload_request_id, file_chunk_directory, file_name, file_size, number_of_chunks) VALUES (?, ?, ?, ?, ?)" +
		"ON DUPLICATE KEY UPDATE file_chunk_directory=VALUES(file_chunk_directory), file_name=VALUES(file_name), file_size=VALUES(file_size), number_of_chunks=VALUES(number_of_chunks)"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, uploadRequestId, fileChunkDir, fileName, fileSize, numberOfChunks)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// ON DUPLICATE KEY UPDATE will result in 2 rows being updated if there is a duplicate key
	if rows == 0 {
		log.Warnf("Connection request information already exists in database for connection request with ID: %s", uploadRequestId)
	} else if rows == 1 {
		log.Infof("Connection request inserted into database for connection request with ID: %s", uploadRequestId)
	} else {
		log.Info("Connection request updated in database for connection request with ID: %s", uploadRequestId)
	}

	return nil
}
