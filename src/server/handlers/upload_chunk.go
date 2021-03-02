package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/jobs"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"io"
	"os"
	"path"
)

var (
	uploadCounter = 0
)

func reassembleUponAllChunksReceived(cache data.Cache, db data.Database, numberOfChunks int, uploadRequestId string, directoryId string) error {
	// Check if upload counter has
	if numberOfChunks == uploadCounter {
		uploadCounter = 0
		fileName, err := data.Database.ShowConnectionRequestFileName(db, uploadRequestId)
		if err != nil {
			return err
		}

		if err = jobs.ReassembleFile(cache, directoryId, fileName, uploadRequestId); err != nil {
			log.Fatalf("HEEEEEEEEEEEERE %v", err)
			return err
		}
	}

	return nil
}

func createChunkFile(filename string, chunkData io.Reader) middleware.Responder {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return middleware.Error(500, fmt.Errorf("could not create file on server %v", err))
	}

	n, err := io.Copy(f, chunkData)
	if err != nil {
		return middleware.Error(500, fmt.Errorf("could not upload file on server %v", err))
	}

	log.Infof("Copied bytes %d", n)
	log.Infof("File uploaded copied as %s", filename)

	return nil
}

func UploadChunkHandler(params files.UploadChunkParams, dataAccess data.Wrapper) middleware.Responder {
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
	db := dataAccess.Database
	uploadRequestId := params.UploadRequestID
	chunkHash := params.ChunkHash

	// TODO: Check here that the actual hash of the data matches the provided hash (check params.ChunkHash == real hash)

	// uploads file and save it locally
	directory := path.Join(c.Config.Server.UploadDir, uploadRequestId)
	filename := path.Join(directory, fmt.Sprintf("%d_%s", params.ChunkNumber, params.ChunkHash))

	uploadCounter++

	// Create the file to write the chunk to
	if err := createChunkFile(filename, params.ChunkData); err != nil {
		return middleware.Error(409, fmt.Errorf("could not access file data on server: %v", err))
	}

	// Get directory ID from connection_request table
	directoryId, err := data.Database.ShowFileChunkDirectory(db, uploadRequestId)
	if err != nil {
		return middleware.Error(400, fmt.Errorf("could not access file data on server: %v", err))
	}

	// Insert chunk info into data
	if err := data.Database.InsertChunk(db, namedFile.Header.Filename, namedFile.Header.Size, chunkHash, params.ChunkNumber, uploadRequestId, directoryId); err != nil {
		return middleware.Error(400, fmt.Errorf("could not store file info on server: %v", err))
	}

	// Reassemble file when uploadCounter indicates all chunks have been received
	numberOfChunks, err := dataAccess.RetrieveNumberOfChunks(uploadRequestId)
	if err != nil {
		return middleware.Error(501, fmt.Errorf("could not access file data on server: %v", err))
	}

	if err := reassembleUponAllChunksReceived(dataAccess.Cache, db, int(numberOfChunks), uploadRequestId, directoryId); err != nil {
		return middleware.Error(400, fmt.Errorf("could not reassemble file on server: %v", err))
	}

	return files.NewUploadChunkCreated().WithPayload(&models.UploadedChunk{
		CompositeFileID: uploadRequestId,
		DirectoryID:     models.DirectoryID(directoryId),
		Hash:            chunkHash,
	})
}
