package data

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

type Chunk struct {
	fileName   string
	fileSize    string
	chunkHash   string
	chunkNumber int64
	uploadRequestId    string
	fileChunkDirectory string
}

type ConnectionRequest struct {
	uploadRequestId    string
	fileChunkDirectory string
	fileName         string
	fileSize       int64
	numberOfChunks int64
}

func (db *DatabaseData) InsertChunk(fileName string, fileSize int64, chunkHash string, chunkNumber int64, uploadRequestId string, directoryId string) error {
	query := "INSERT INTO chunk(file_name, file_size, chunk_hash, chunk_number, upload_request_id, file_chunk_directory) VALUES (?, ?, ?, ?, ?, ?)" +
		"ON DUPLICATE KEY UPDATE file_name=VALUES(file_name), file_size=VALUES(file_size), chunk_number=VALUES(chunk_number), upload_request_id=VALUES(upload_request_id), file_chunk_directory=VALUES(file_chunk_directory)"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	stmt, err := db.mysql.PrepareContext(ctx, query)
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
		log.Warnf("Chunk information already exists in db for chunk: %s", chunkHash)
	} else if rows == 1 {
		log.Infof("Chunk information inserted into db for chunk: %s", chunkHash)
	} else {
		log.Infof("Chunk information updated in db for chunk: %s", chunkHash)
	}

	return nil
}

func (db *DatabaseData) InsertConnectionRequest(uploadRequestId string, fileChunkDir string, fileName string, fileSize int64, numberOfChunks int64) error {
	query := "INSERT INTO connection_request(upload_request_id, file_chunk_directory, file_name, file_size, number_of_chunks) VALUES (?, ?, ?, ?, ?)" +
		"ON DUPLICATE KEY UPDATE file_chunk_directory=VALUES(file_chunk_directory), file_name=VALUES(file_name), file_size=VALUES(file_size), number_of_chunks=VALUES(number_of_chunks)"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	stmt, err := db.mysql.PrepareContext(ctx, query)
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
		log.Warnf("Connection request information already exists in db for connection request with ID: %s", uploadRequestId)
	} else if rows == 1 {
		log.Infof("Connection request inserted into db for connection request with ID: %s", uploadRequestId)
	} else {
		log.Infof("Connection request updated in db for connection request with ID: %s", uploadRequestId)
	}

	return nil
}

func (db *DatabaseData) ShowNumberOfChunks(uploadRequestId string) (numberOfChunks int64, err error) {
	query := "SELECT number_of_chunks FROM connection_request WHERE upload_request_id=?"
	row := db.mysql.QueryRow(query, uploadRequestId)

	var connectionRequest ConnectionRequest
	err = row.Scan(&connectionRequest.numberOfChunks)

	return connectionRequest.numberOfChunks, err
}

func (db *DatabaseData) ShowConnectionRequestFileName(uploadRequestId string) (fileName string, err error) {
	query := "SELECT file_name FROM connection_request WHERE upload_request_id=?"
	row := db.mysql.QueryRow(query, uploadRequestId)

	var connectionRequest ConnectionRequest
	err = row.Scan(&connectionRequest.fileName)

	return connectionRequest.fileName, err
}

func (db *DatabaseData) ShowFileChunkDirectory(uploadRequestId string) (directoryId string, err error) {
	query := "SELECT file_chunk_directory FROM connection_request WHERE upload_request_id=?"
	col := db.mysql.QueryRow(query, uploadRequestId)

	var connectionRequest ConnectionRequest
	err = col.Scan(&connectionRequest.fileChunkDirectory)

	return connectionRequest.fileChunkDirectory, err
}
