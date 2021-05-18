package repo

import (
	"errors"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gorm.io/gorm"
	"strconv"
)

func GetFileChunksInOrder(fileID uint, db *gorm.DB) (chunks []schema.FileChunk, err error) {
	tx := db.Joins("Chunk").Where("file_id = ?", fileID).Order("number ASC").Find(&chunks)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return chunks, nil
}

func GetCachedChunksReceived(fileID uint64) (uint64, error) {
	item, found := FileIDChunkCountCache.Get(strconv.Itoa(int(fileID)))
	if found {
		chunksReceived, ok := item.(uint64)
		if ok {
			return chunksReceived, nil
		}
	}
	return 0, errors.New("invalid cache entry for chunks received")
}
