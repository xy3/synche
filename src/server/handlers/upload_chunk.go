package handlers

import (
	"database/sql"
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"io"
	"log"
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
		log.Print("=== Received new chunk ===")
		log.Printf("Filename: %s", namedFile.Header.Filename)
		log.Printf("Size: %d", namedFile.Header.Size)
		log.Printf("ChunkHash: %s", params.ChunkHash)
		log.Printf("ChunkNumber: %d", params.ChunkNumber)
		log.Printf("UploadRequestID: %s", params.UploadRequestID)
	}

	// Frequently used params
	uploadRequestId := params.UploadRequestID
	chunkHash := params.ChunkHash

	// TODO: Check here that the actual hash of the data matches the provided hash (check params.ChunkHash == real hash)

	// uploads file and save it locally
	directory := path.Join(viper.GetString("server.uploadDirectory"), uploadRequestId)
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

	log.Printf("Copied bytes %d", n)

	log.Printf("File uploaded copied as %s", filename)

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
		log.Printf("Could not insert chunk %s", params.ChunkHash)
	} else {
		log.Printf("Inserted chunk information into database for chunk: %s", params.ChunkHash)
	}

	return files.NewUploadChunkCreated().WithPayload(&models.UploadedChunk{
		CompositeFileID: uploadRequestId,
		DirectoryID:     models.DirectoryID(directoryId),
		Hash:            chunkHash,
	})
}