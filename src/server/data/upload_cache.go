package data

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
)

const uploadPrefix = "upload"

//go:generate mockery --name=UploadCache --case underscore
type UploadCache interface {
	FormatKey(uploadId uint) string
	GetUpload(cache Cache, uploadId uint) (upload schema.Upload, err error)
	SetUpload(cache Cache, uploadId uint, upload schema.Upload) error
	DelUpload(cache Cache, uploadId uint) error
}

type uploadCache struct {}

func NewUploadCache() *uploadCache {
	return &uploadCache{}
}

func (c *uploadCache) FormatKey(uploadId uint) string {
	return fmt.Sprintf("%s:%d", uploadPrefix, uploadId)
}

func (c *uploadCache) GetUpload(cache Cache, uploadId uint) (upload schema.Upload, err error) {
	res, err := redis.Values(cache.GetAll(c.FormatKey(uploadId)))
	if err != nil {
		return upload, err
	}
	err = redis.ScanStruct(res, &upload)
	return upload, err
}

func (c *uploadCache) SetUpload(cache Cache, uploadId uint, upload schema.Upload) error {
	_, err := cache.SetAll(c.FormatKey(uploadId), upload)
	return err
}

func (c *uploadCache) DelUpload(cache Cache, uploadId uint) error {
	_, err := cache.Delete(c.FormatKey(uploadId))
	return err
}
