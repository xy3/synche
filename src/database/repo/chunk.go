package repo

import (
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"strconv"
)

func GetFileChunksInOrder(fileID uint) ([]schema.FileChunk, error) {
	var chunks []schema.FileChunk
	tx := database.DB.Joins("Chunk").Where("file_id = ?", fileID).Order("number ASC").Find(&chunks)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return chunks, nil
}

func GetCachedChunksReceived(fileID uint64) uint64 {
	item, found := FileIDChunkCountCache.Get(strconv.Itoa(int(fileID)))
	if found {
		chunksReceived, ok := item.(uint64)
		if ok {
			return chunksReceived
		}
	}
	log.Error("invalid cache entry for chunks received")
	return 0
}
