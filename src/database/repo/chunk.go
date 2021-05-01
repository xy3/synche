package repo

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
)

func GetFileChunksInOrder(fileID uint) ([]schema.FileChunk, error) {
	var chunks []schema.FileChunk
	tx := database.DB.Joins("Chunk").Where("file_id = ?", fileID).Order("number ASC").Find(&chunks)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return chunks, nil
}