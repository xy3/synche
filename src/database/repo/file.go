package repo

import (
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"io"
	"path/filepath"
)

func GetFileByID(fileID uint) (*schema.File, error) {
	var file schema.File
	if err := database.DB.First(&file, fileID).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// FindFileByFullPath finds a file by its full path on the disk. It assumes the path given
// is a file and not a directory, so if it is given the path to a directory it will treat it
// like a path to a file.
func FindFileByFullPath(path string) (*schema.File, error) {
	if fileValue, ok := pathFileCache.Get(path); ok {
		return fileValue.(*schema.File), nil
	}

	// Get the md5 hash of only the directory part of the path
	dirPath := filepath.Dir(path)
	dirPathHash := hash.PathHash(dirPath)
	log.Infof("Received FindFileByFullPath request for file: %s", path)
	log.Infof("Received FindFileByFullPath request with dir as: %s", dirPath)

	file := &schema.File{}
	tx := database.DB.Model(schema.File{}).Joins("Directory").Where(schema.File{
		Directory: &schema.Directory{PathHash: dirPathHash},
	}).First(file)

	if tx.Error != nil {
		return nil, tx.Error
	}

	log.Info(file)
	_ = pathFileCache.Add(path, file, cache.DefaultExpiration)
	return file, nil
}

func CreateFileFromReader(path string, reader io.Reader, num int, rootChunkPath *string) (
	*schema.File,
	error,
) {
	// TODO
	return nil, nil
}


func RenameFile(fileID uint, newName string) error {
	tx := database.DB.Model(schema.File{}).Where("id = ?", fileID).Update("name", newName)
	return tx.Error
}