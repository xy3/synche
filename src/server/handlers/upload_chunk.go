package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/jobs"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
	"io"
	"path/filepath"
)

var (
	uploadCounter = 0
)

func reassembleFile(db data.SyncheData, uploadRequestId string, directoryId string) error {
	fileName, err := db.Database.ConnectionRequestFileName(uploadRequestId)
	if err != nil {
		return err
	}
	return jobs.ReassembleFile(db.Cache, directoryId, fileName, uploadRequestId)
}

func writeChunkFile(filename string, chunkData io.Reader) error {
	return afero.WriteReader(files.AppFS, filename, chunkData)
}

func UploadChunkHandler(params transfer.UploadChunkParams, data data.SyncheData) middleware.Responder {
	if params.ChunkData == nil {
		return transfer.NewUploadChunkBadRequest().WithPayload("no chunk data received")
	}
	defer params.ChunkData.Close()

	namedFile, ok := params.ChunkData.(*runtime.File)
	if ok {
		log.WithFields(log.Fields{
			"filename":          namedFile.Header.Filename,
			"size":              namedFile.Header.Size,
			"chunk hash":        params.ChunkHash,
			"chunk number":      params.ChunkNumber,
			"upload request id": params.UploadRequestID,
		}).Info("Received new chunk")
	}

	// TODO: Check here that the actual hash of the data matches the provided hash (check params.ChunkHash == real hash)

	directory := filepath.Join(c.Config.Server.UploadDir, params.UploadRequestID)
	chunkFilename := filepath.Join(directory, fmt.Sprintf("%d_%s", params.ChunkNumber, params.ChunkHash))

	uploadCounter++

	err := writeChunkFile(chunkFilename, params.ChunkData)
	if err != nil {
		return transfer.NewUploadChunkConflict().WithPayload(models.Error(err.Error()))
	}
	return storeChunkData(data, namedFile, params)

}

func storeChunkData(
	data data.SyncheData,
	file *runtime.File,
	params transfer.UploadChunkParams,
) middleware.Responder {
	// Get directory ID from connection_request table
	directoryId, err := data.Database.ChunkDirectory(params.UploadRequestID)
	if err != nil {
		return transfer.NewUploadChunkBadRequest().WithPayload("failed to find the chunk directory")
	}

	// Insert chunk info into data
	err = data.Database.InsertChunk(file.Header.Filename, file.Header.Size, params.ChunkHash, params.ChunkNumber, params.UploadRequestID, directoryId)
	if err != nil {
		return transfer.NewUploadChunkBadRequest().WithPayload("failed to add the chunk data to the database")
	}

	// Reassemble file when uploadCounter indicates all chunks have been received
	numberOfChunks, err := data.NumberOfChunks(params.UploadRequestID)
	if err != nil {
		return transfer.NewUploadChunkBadRequest().WithPayload("failed to find the number of chunks for this file")
	}

	if uploadCounter >= int(numberOfChunks) {
		uploadCounter = 0
		err = reassembleFile(data, params.UploadRequestID, directoryId)
		if err != nil {
			return transfer.NewUploadChunkBadRequest().WithPayload("failed to re-assemble the file")
		}
	}

	return transfer.NewUploadChunkCreated().WithPayload(&models.UploadedChunk{
		CompositeFileID: params.UploadRequestID,
		DirectoryID:     models.DirectoryID(directoryId),
		Hash:            params.ChunkHash,
	})
}
