package handlers

import (
	"database/sql"
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"io"
	"os"
	"path"
)

var (
	uploadCounter = 0
)

func UploadChunkHandler(params files.UploadChunkParams, db *sql.DB) middleware.Responder {
	if params.ChunkData == nil {
		return middleware.Error(404, fmt.Errorf("no file provided"))
	}
	defer params.ChunkData.Close()

	namedFile, ok := params.ChunkData.(*runtime.File)
	if ok {
		log.Info("=== Received new chunk ===")
		log.Infof("Filename: %s", namedFile.Header.Filename)
		log.Infof("Size: %d", namedFile.Header.Size)
		log.Infof("ChunkHash: %s", params.ChunkHash)
		log.Infof("ChunkNumber: %d", params.ChunkNumber)
		log.Infof("UploadRequestID: %s", params.UploadRequestID)
	}

	// Frequently used params
	uploadRequestId := params.UploadRequestID
	chunkHash := params.ChunkHash

	// TODO: Check here that the actual hash of the data matches the provided hash (check params.ChunkHash == real hash)

	// uploads file and save it locally
	directory := path.Join(c.Config.Server.UploadDir, uploadRequestId)
	filename := path.Join(directory, fmt.Sprintf("%d_%s", params.ChunkNumber, params.ChunkHash))

	uploadCounter++

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return middleware.Error(500, fmt.Errorf("could not create file on server"))
	}

	n, err := io.Copy(f, params.ChunkData)
	if err != nil {
		return middleware.Error(500, fmt.Errorf("could not upload file on server"))
	}

	log.Infof("Copied bytes %d", n)

	log.Infof("File uploaded copied as %s", filename)

	// Get directory ID from connection_request table
	var directoryId string
	query := "SELECT file_chunk_directory FROM connection_request WHERE upload_request_id=?"
	col := db.QueryRow(query, uploadRequestId)
	if err := col.Scan(&directoryId); err != nil {
		log.Fatal("Could not find row directory_id in table connection_request")
	}

	// Insert chunk info into database
	err = database.InsertChunk(db, namedFile.Header.Filename, namedFile.Header.Size, chunkHash, params.ChunkNumber, uploadRequestId, directoryId)
	if err != nil {
		log.Fatalf("Could not insert chunk %s\n-----> %s", params.ChunkHash, err)
	}

	return files.NewUploadChunkCreated().WithPayload(&models.UploadedChunk{
		CompositeFileID: uploadRequestId,
		DirectoryID:     models.DirectoryID(directoryId),
		Hash:            chunkHash,
	})
}
