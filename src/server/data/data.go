package data

import (
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

type SyncheData struct {
	Cache    *Cache
	Database *Database
}

func (d *SyncheData) NumberOfChunks(uploadRequestId string) (chunks int64, err error) {
	// Check cache first
	chunks, err = d.Cache.GetNumberOfChunks(uploadRequestId)

	// Get data from db if it wasn't in the cache
	if err == redis.ErrNil {
		log.Warnf("Upload request with ID: %v was not retreived from cache", uploadRequestId)
		chunks, err = d.Database.NumberOfChunks(uploadRequestId)
		if err != nil {
			return 0, err
		}

		// Add data to cache for future usage
		err = d.Cache.SetNumberOfChunks(uploadRequestId, chunks)
		if err != nil {
			return 0, err
		}

		return chunks, nil
	}

	return
}
