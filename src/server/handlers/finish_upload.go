package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
)

func getMissingChunks(fileID uint64, expectedNumOfChunks int64) ([]int64, error) {
	missingChunks := make([]int64, 0)

	// get all received chunks
	fileChunks, err := repo.GetFileChunksInOrder(uint(fileID))
	if err != nil {
		return nil, err
	}

	// iterate through received chunks to see if any chunk is missing
	for i := int64(0); i < expectedNumOfChunks; i++ {
		if fileChunks[i].Number != i+1 {
			missingChunks = append(missingChunks, i+1)
		}
	}

	return missingChunks, nil
}

func FinishUpload(params transfer.FinishUploadParams, user *schema.User) middleware.Responder {
	fileID := params.FileID
	chunksReceived := repo.GetCachedChunksReceived(fileID)
	expectedNumOfChunks, err := repo.GetTotalFileChunks(fileID)
	if err != nil || chunksReceived == 0 {
		return transfer.NewFinishUploadDefault(500).WithPayload("failed to access cache")
	}

	// If the amount of chunks received != the expected, find what chunks are missing
	if chunksReceived != expectedNumOfChunks {
		missingChunks, err := getMissingChunks(fileID, int64(expectedNumOfChunks))
		if err != nil {
			return transfer.NewFinishUploadDefault(500).WithPayload("failed to access database")
		}
		return transfer.NewFinishUploadOK().WithPayload(&models.MissingChunks{ChunkNumbers: missingChunks})
	}

	// No missing chunks
	return transfer.NewFinishUploadOK().WithPayload(&models.MissingChunks{ChunkNumbers: []int64{}})
}
