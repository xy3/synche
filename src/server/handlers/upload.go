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
	UploadDirectory = "../data/received/"
)

var (
	uploadCounter = 0
)


func UploadFileHandler(params files.UploadFileParams) middleware.Responder {
	if params.ChunkData == nil {
		return middleware.Error(404, fmt.Errorf("no file provided"))
	}
	defer params.ChunkData.Close()

	if namedFile, ok := params.ChunkData.(*runtime.File); ok {
		log.Printf("received chunk name: %s", namedFile.Header.Filename)
		log.Printf("received chunk size: %d", namedFile.Header.Size)
		log.Printf("received chunk hash: %d", params.ChunkHash)
		log.Printf("received chunk number: %d", params.ChunkNumber)
		log.Printf("received chunk request id: %d", params.UploadRequestID)
	}

	// TODO: Check here that the actual hash of the data matches the provided hash (check params.ChunkHash == real hash)

	// uploads file and save it locally
	filename := path.Join(UploadDirectory, fmt.Sprintf("%s_%d_%d", params.ChunkHash, params.ChunkNumber, uploadCounter))
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

	return files.NewUploadFileCreated().WithPayload(&models.UploadedChunk{
		CompositeFileID: params.UploadRequestID,
		DirectoryID:     UploadDirectory, // TODO: Change this to the UUID of the directory from the database
		Hash:            params.ChunkHash,
	})
}