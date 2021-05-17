package handlers

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/jobs"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
	"path/filepath"
	"strconv"
)

var (
	badRequest   = transfer.NewUploadChunkDefault(400)
	serverErr    = transfer.NewUploadChunkDefault(500)
	fileConflict = transfer.NewUploadChunkDefault(409)
	errNoData    = badRequest.WithPayload("no chunk data received")
)

// UploadChunk Handles new chunks being uploaded and responds to the client with each chunk status
func UploadChunk(params transfer.UploadChunkParams, user *schema.User) middleware.Responder {
	if params.ChunkData == nil {
		return errNoData
	}
	defer params.ChunkData.Close()

	namedFile, ok := params.ChunkData.(*runtime.File)
	if ok {
		log.WithFields(log.Fields{
			"Size":        namedFile.Header.Size,
			"ChunkHash":   params.ChunkHash,
			"ChunkNumber": params.ChunkNumber,
		}).Info("Received new chunk")
	}

	chunkBytes, err := afero.ReadAll(params.ChunkData)
	if err != nil {
		return badRequest.WithPayload("Failed to read the chunk bytes")
	}

	if !hash.ValidateChunkHash(params.ChunkHash, chunkBytes) {
		return badRequest.WithPayload("chunk hash does not match its data")
	}

	file, err := repo.GetFileByID(uint(params.FileID), database.DB)
	if err != nil {
		return badRequest.WithPayload("Failed to find a related file")
	}

	if err = writeChunkFile(chunkBytes, c.Config.Server.ChunkDir, params.ChunkHash); err != nil {
		return fileConflict.WithPayload(models.Error(err.Error()))
	}

	return storeChunkData(namedFile, params, file)
}

// writeChunkFile Creates and writes the chunk details to the chunk file
func writeChunkFile(chunkData []byte, chunkDir, chunkHash string) error {
	chunkFilename := filepath.Join(chunkDir, chunkHash)
	return files.Afs.WriteFile(chunkFilename, chunkData, 0644)
}

// storeChunkData Adds chunk details to the database and cache
func storeChunkData(
	chunkFile *runtime.File,
	params transfer.UploadChunkParams,
	file *schema.File,
) middleware.Responder {

	db := database.DB.Begin()

	chunk := schema.Chunk{
		Hash: params.ChunkHash,
		Size: chunkFile.Header.Size,
	}

	if db.Where(chunk).FirstOrCreate(&chunk).Error != nil {
		db.Rollback()
		return serverErr.WithPayload("failed to add the chunk data to the database")
	}

	// Insert chunk info into data
	fileChunk := schema.FileChunk{
		Number:  params.ChunkNumber,
		ChunkID: chunk.ID,
		FileID:  file.ID,
	}

	if db.Where(fileChunk).FirstOrCreate(&fileChunk).Error != nil {
		db.Rollback()
		return serverErr.WithPayload("failed to add the chunk data to the database")
	}

	db.Commit()

	// err := file.Upload.UpdateChunksReceived(1, database.DB)
	// if err != nil {
	// 	return serverErr.WithPayload("failed to update the number of chunks received")
	// }
	//

	strFileID := strconv.Itoa(int(file.ID))
	var chunksReceived int64
	if item, ok := repo.FileIDChunkCountCache.Get(strFileID); ok {
		chunksReceived, ok = item.(int64)
		if !ok {
			log.Error("invalid cache entry for chunks received")
		}
	} else {
		repo.FileIDChunkCountCache.Set(strFileID, int64(0), cache.NoExpiration)
	}

	err := repo.FileIDChunkCountCache.Increment(strFileID, 1)
	if err != nil {
		return serverErr.WithPayload("failed increment the chunks received count")
	}
	chunksReceived += 1

	log.Infof("Received chunks: %d", chunksReceived)

	if chunksReceived >= file.TotalChunks {
		chunksReceived = 0
		// delete(chunksReceived, file.ID)
		repo.FileIDChunkCountCache.Delete(strFileID)
		log.Infof("Reassembling the file: %d", file.ID)
		err = jobs.ReassembleFile(c.Config.Server.ChunkDir, file)
		if err != nil {
			return badRequest.WithPayload("Failed to re-assemble the file")
		}
	}

	return transfer.NewUploadChunkCreated().WithPayload(&models.FileChunk{
		Chunk: &models.Chunk{
			Hash: chunk.Hash,
			ID:   uint64(chunk.ID),
			Size: chunk.Size,
		},
		FileID: uint64(file.ID),
		ID:     uint64(fileChunk.ID),
		Number: fileChunk.Number,
	})
}
