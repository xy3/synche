package data

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

const (
	connectionRequestPrefix string = "connection_request:"
)

type ConnectionRequestCache struct {
	UploadRequestId string `json:"upload_request_id"`
	NumberOfChunks  int64  `json:"number_of_chunks,string"`
}

func (c *CacheData) SetNumberOfChunks(uploadRequestId string, numberOfChunks int64) error {
	// Create object to add to cache
	crc := ConnectionRequestCache{ UploadRequestId: uploadRequestId,
								   NumberOfChunks: numberOfChunks}

	result, err := json.Marshal(crc)
	if err != nil {
		return err
	}

	// Get a connection from the pool
	conn := c.redis.Get()

	// Add object to cache
	_, err = conn.Do("SET", connectionRequestPrefix+ crc.UploadRequestId, result)
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheData) GetNumberOfChunks(uploadRequestId string) (numberOfChunks int64, err error) {
	// Get a connection from the pool
	conn := c.redis.Get()

	// Get object from cache
	res, err := redis.String(conn.Do("Get", connectionRequestPrefix + uploadRequestId))
	if err == redis.ErrNil {
		log.Errorf("Upload request not cached: %s", uploadRequestId)
		return 0, err
	} else if err != nil {
		return 0, err
	}
	crc := ConnectionRequestCache{}
	err = json.Unmarshal([]byte(res), &crc)

	return crc.NumberOfChunks, err
}

func (c *CacheData) DeleteKeyNumberOfChunks(uploadRequestId string) error {
	// Get a connection from the pool
	conn := c.redis.Get()

	_, err := redis.Int(conn.Do("DEL", connectionRequestPrefix + uploadRequestId))
	if err != nil {
		return err
	}

	return nil
}
