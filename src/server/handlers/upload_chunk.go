package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"io"
	"log"
	"os"
	"path"
)

const (
	UploadDirectory = "../data/received/" // TODO: get this value from a stored config file
)

var (
	uploadCounter = 0
)


func UploadChunkHandler(params files.UploadChunkParams) middleware.Responder {
	if params.ChunkData == nil {
		return middleware.Error(404, fmt.Errorf("no file provided"))
	}
	defer params.ChunkData.Close()

	if namedFile, ok := params.ChunkData.(*runtime.File); ok {
		log.Print("=== Received new chunk ===")
		log.Printf("Filename: %s", namedFile.Header.Filename)
		log.Printf("Size: %d", namedFile.Header.Size)
		log.Printf("ChunkHash: %s", params.ChunkHash)
		log.Printf("ChunkNumber: %d", params.ChunkNumber)
		log.Printf("UploadRequestID: %s", params.UploadRequestID)
	}

	// TODO: Check here that the actual hash of the data matches the provided hash (check params.ChunkHash == real hash)

	// uploads file and save it locally
	directory := path.Join(UploadDirectory, params.UploadRequestID)
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

	log.Printf("copied bytes %d", n)

	log.Printf("file uploaded copied as %s", filename)

	return files.NewUploadChunkCreated().WithPayload(&models.UploadedChunk{
		CompositeFileID: params.UploadRequestID,
		DirectoryID:     models.DirectoryID(directory), // TODO: Change this to the ID of the directory from the database
		Hash:            params.ChunkHash,
	})
}