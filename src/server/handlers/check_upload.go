package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
)

// CheckUpload checks what chunks have already been received for a given upload
func CheckUpload(params transfer.CheckUploadedChunksParams, user *schema.User) middleware.Responder {
	fileChunks, err := repo.GetFileChunksInOrder(uint(params.FileID), database.DB)
	if err != nil {
		return transfer.NewCheckUploadedChunksDefault(500).WithPayload("failed to find existing chunks")
	}

	var existingChunks []int64
	for _, chunk := range fileChunks {
		existingChunks = append(existingChunks, chunk.Number)
	}

	return transfer.NewCheckUploadedChunksOK().WithPayload(&models.ExistingChunks{ChunkNumbers: existingChunks})
}
