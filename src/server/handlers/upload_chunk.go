package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/jobs"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
	"path/filepath"
)

var (
	badRequest    = transfer.NewNewUploadBadRequest()
	fileConflict  = transfer.NewUploadChunkConflict()
	uploadCounter = 0
)

func UploadChunkHandler(params transfer.UploadChunkParams, syncheData *data.SyncheData) middleware.Responder {
	if params.ChunkData == nil {
		return transfer.NewUploadChunkBadRequest().WithPayload("no chunk data received")
	}
	defer params.ChunkData.Close()

	namedFile, ok := params.ChunkData.(*runtime.File)
	if ok {
		log.WithFields(log.Fields{
			"Size":            namedFile.Header.Size,
			"ChunkHash":       params.ChunkHash,
			"ChunkNumber":     params.ChunkNumber,
			"UploadRequestID": params.UploadRequestID,
		}).Info("Received new chunk")
	}

	chunkBytes, err := afero.ReadAll(params.ChunkData)
	if err != nil {
		return badRequest.WithPayload("failed to read the chunk bytes")
	}

	if !files.ValidateChunkHash(params.ChunkHash, chunkBytes) {
		return badRequest.WithPayload("chunk hash does not match its data")
	}

	// Todo: store this in the database (or redis?) as 'chunksUploaded'
	uploadCounter++

	var upload schema.Upload
	tx := syncheData.DB.Joins("Directory").Joins("File")
	if tx.First(&upload, params.UploadRequestID).Error != nil {
		transfer.NewNewUploadBadRequest().WithPayload("failed to find a related upload request")
	}

	if err = writeChunkFile(chunkBytes, params.ChunkNumber, upload.Directory.Path, params.ChunkHash); err != nil {
		return fileConflict.WithPayload(models.Error(err.Error()))
	}

	return storeChunkData(syncheData, namedFile, params, upload)
}

func writeChunkFile(chunkData []byte, chunkNumber int64, chunkDir, chunkHash string) error {
	chunkFilename := filepath.Join(chunkDir, fmt.Sprintf("%d_%s", chunkNumber, chunkHash))
	return files.Afs.WriteFile(chunkFilename, chunkData, 0644)
}

func storeChunkData(
	syncheData *data.SyncheData,
	chunkFile *runtime.File,
	params transfer.UploadChunkParams,
	upload schema.Upload,
) middleware.Responder {
	// Insert chunk info into data
	newChunk := schema.FileChunk{
		Chunk: schema.Chunk{
			Hash: params.ChunkHash,
			Size: chunkFile.Header.Size,
		},
		Number:      params.ChunkNumber,
		DirectoryID: upload.DirectoryID,
		FileID:      upload.FileID,
		UploadID:    upload.ID,
	}

	tx := syncheData.DB.Begin()

	if tx.Create(&newChunk).Error != nil {
		tx.Rollback()
		return badRequest.WithPayload("failed to add the chunk data to the database")
	}

	tx.Commit()

	// Reassemble file when uploadCounter indicates all chunks have been received
	if uploadCounter >= int(upload.NumChunks) {
		uploadCounter = 0
		err := jobs.ReassembleFile(syncheData.Cache, upload.Directory.Path, upload.File.Name, upload.ID)
		if err != nil {
			return badRequest.WithPayload("failed to re-assemble the file")
		}
	}

	return transfer.NewUploadChunkCreated().WithPayload(&models.UploadedChunk{
		DirectoryID: int64(upload.DirectoryID),
		FileID:      int64(upload.FileID),
		Hash:        params.ChunkHash,
	})
}
