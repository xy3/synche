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
	badRequest    = transfer.NewNewUploadDefault(400)
	fileConflict  = transfer.NewUploadChunkConflict()
	uploadCounter = 0
)

func UploadChunk(
	params transfer.UploadChunkParams,
	user *schema.User,
) middleware.Responder {
	if params.ChunkData == nil {
		return transfer.NewUploadChunkBadRequest().WithPayload("no chunk data received")
	}
	defer params.ChunkData.Close()

	namedFile, ok := params.ChunkData.(*runtime.File)
	if ok {
		log.WithFields(log.Fields{
			"Size":        namedFile.Header.Size,
			"ChunkHash":   params.ChunkHash,
			"ChunkNumber": params.ChunkNumber,
			"UploadID":    params.UploadID,
		}).Info("Received new chunk")
	}

	chunkBytes, err := afero.ReadAll(params.ChunkData)
	if err != nil {
		return badRequest.WithPayload("Failed to read the chunk bytes")
	}

	if !files.ValidateChunkHash(params.ChunkHash, chunkBytes) {
		return badRequest.WithPayload("chunk hash does not match its data")
	}

	// Todo: store this in the database (or redis?) as 'chunksUploaded'
	uploadCounter++

	var upload schema.Upload
	tx := data.DB.Joins("ChunkDirectory").Joins("File")
	if tx.First(&upload, params.UploadID).Error != nil {
		badRequest.WithPayload("Failed to find a related upload request")
	}

	if err = writeChunkFile(chunkBytes, params.ChunkNumber, upload.ChunkDirectory.Path, params.ChunkHash); err != nil {
		return fileConflict.WithPayload(models.Error(err.Error()))
	}

	return storeChunkData(namedFile, params, upload)
}

func writeChunkFile(chunkData []byte, chunkNumber int64, chunkDir, chunkHash string) error {
	chunkFilename := filepath.Join(chunkDir, fmt.Sprintf("%d_%s", chunkNumber, chunkHash))
	return files.Afs.WriteFile(chunkFilename, chunkData, 0644)
}

func storeChunkData(
	chunkFile *runtime.File,
	params transfer.UploadChunkParams,
	upload schema.Upload,
) middleware.Responder {
	// Insert chunk info into data
	fileChunk := schema.FileChunk{
		Number: params.ChunkNumber,
		Chunk: schema.Chunk{
			Hash: params.ChunkHash,
			Size: chunkFile.Header.Size,
		},
		ChunkDirectoryID: upload.ChunkDirectoryID,
		ChunkDirectory:   schema.Directory{},
		FileID:           upload.FileID,
		File:             schema.File{},
		UploadID:         upload.ID,
		Upload:           schema.Upload{},
	}

	tx := data.DB.Begin()

	if tx.Create(&fileChunk).Error != nil {
		tx.Rollback()
		return badRequest.WithPayload("Failed to add the chunk data to the database")
	}

	tx.Commit()

	// Reassemble file when uploadCounter indicates all chunks have been received
	if uploadCounter >= int(upload.NumChunks) {
		uploadCounter = 0

		err := jobs.ReassembleFile(upload.ChunkDirectory.Path, upload.File.Name, upload.ID)

		if err != nil {
			return badRequest.WithPayload("Failed to re-assemble the file")
		}
	}

	return transfer.NewUploadChunkCreated().WithPayload(&models.FileChunk{
		Chunk: &models.Chunk{
			Hash: fileChunk.Chunk.Hash,
			ID:   uint64(fileChunk.ID),
			Size: fileChunk.Chunk.Size,
		},
		ChunkDirectoryID: uint64(upload.ChunkDirectoryID),
		FileID:           uint64(upload.FileID),
		ID:               uint64(fileChunk.ID),
		Number:           fileChunk.Number,
		UploadID:         uint64(fileChunk.UploadID),
	})
}
