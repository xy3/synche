package data

import (
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

type ServerEnv interface {
	// Methods that may need to access cache and/or database
	RetrieveNumberOfChunks()

	// Pure access
	Cache() Cache
	Database() Database
}

type Wrapper struct {
	Cache    Cache
	Database Database
}

func (w *Wrapper) GetCache() Cache {
	return w.Cache
}

func (w *Wrapper) GetDatabase() Database {
	return w.Database
}

func (w *Wrapper) RetrieveNumberOfChunks(uploadRequestId string) (numberOfChunks int64, err error) {
	// Check cache first
	numberOfChunks, err = w.Cache.GetNumberOfChunks(uploadRequestId)

	// Get data from db if it wasn't in the cache
	if err == redis.ErrNil {
		log.Warnf("Upload request with ID: %v was not retreived from cache", uploadRequestId)
		numberOfChunks, err := w.Database.ShowNumberOfChunks(uploadRequestId)
		if err != nil {
			return 0, err
		}

		// Add data to cache for future usage
		err = w.Cache.SetNumberOfChunks(uploadRequestId, numberOfChunks)
		if err != nil {
			return 0, err
		}

		return numberOfChunks, nil
	}

	if err != nil {
		return 0, err
	}

	return numberOfChunks, nil
}
