package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/xy3/synche/src/server"
	"github.com/xy3/synche/src/server/models"
	"github.com/xy3/synche/src/server/repo"
	"github.com/xy3/synche/src/server/restapi/operations/transfer"
	"github.com/xy3/synche/src/server/schema"
)

// CheckUpload checks what chunks have already been received for a given upload
func CheckUpload(params transfer.CheckUploadedChunksParams, _ *schema.User) middleware.Responder {
	fileChunks, err := repo.GetFileChunksInOrder(uint(params.FileID), server.DB)
	if err != nil {
		return transfer.NewCheckUploadedChunksDefault(500).WithPayload("failed to find existing chunks")
	}

	var existingChunks []int64
	for _, chunk := range fileChunks {
		existingChunks = append(existingChunks, chunk.Number)
	}

	return transfer.NewCheckUploadedChunksOK().WithPayload(&models.ExistingChunks{ChunkNumbers: existingChunks})
}
