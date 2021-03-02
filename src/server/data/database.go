package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
)

type Database interface {
	InsertChunk(fileName string, fileSize int64, chunkHash string, chunkNumber int64, uploadRequestId string, directoryId string)  error
	InsertConnectionRequest(uploadRequestId string, fileChunkDir string, fileName string, fileSize int64, numberOfChunks int64) error
	ShowNumberOfChunks(uploadRequestId string) (numberOfChunks int64, err error)
	ShowConnectionRequestFileName(uploadRequestId string) (fileName string, err error)
	ShowFileChunkDirectory(uploadRequestId string) (directoryId string, err error)
}

type DatabaseData struct {
	// wrap sql connection and any other driver that the data needs
	mysql *sql.DB
}

func BuildDBClient(dbConfig config.DatabaseConfig) *DatabaseData {
	dataSourceName := NewDSN(dbConfig)

	db, err := sql.Open(dbConfig.Driver, dataSourceName + dbConfig.Name)
	if err != nil {
		log.Fatal("Error connecting to the data: ", err.Error())
	}

	return &DatabaseData{mysql: db}
}
