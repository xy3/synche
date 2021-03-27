package data

import (
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gorm.io/gorm"
)

type Wrapper interface {
}

type SyncheData struct {
	Cache *RedisCache
	DB    *gorm.DB
}

func NewSyncheData(cache *RedisCache, db *gorm.DB) (*SyncheData, error) {
	return &SyncheData{Cache: cache, DB: db}, nil
}

func (d *SyncheData) MigrateAll() error {
	return d.DB.AutoMigrate(
		&schema.Chunk{},
		&schema.File{},
		&schema.FileChunk{},
		&schema.Directory{},
		&schema.Upload{},
	)
}

func (d *SyncheData) GetNumChunks(uploadCache UploadCache, uploadId uint) (int64, error) {
	upload, err := uploadCache.GetUpload(d.Cache, uploadId)
	if err == nil {
		return upload.NumChunks, nil
	}

	log.WithFields(log.Fields{"UploadID": uploadId}).Debug("Upload request was not retrieved from cache", uploadId)

	res := d.DB.First(&upload, uploadId)
	if res.Error != nil {
		return 0, res.Error
	}

	err = uploadCache.SetUpload(d.Cache, uploadId, upload)
	if err != nil {
		log.WithFields(log.Fields{"UploadID": uploadId}).Error("Failed to cache the number of chunks")
	}
	return upload.NumChunks, nil
}
