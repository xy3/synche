package repo

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
)

func GetStorageDirContents(dirId uint64) ([]schema.File, error) {
	var files []schema.File
	tx := data.DB.Find(&files, schema.File{StorageDirectoryID: uint(dirId)})
	if tx.Error != nil {
		return files, tx.Error
	}
	return files, nil
}

func GetStorageDirIdByPath(dirPath string) (uint, error) {
	var storageDirectory schema.Directory
	tx := data.DB.Begin()
	if res := tx.First(&storageDirectory, "path = ?", dirPath); res.Error != nil {
		return 0, res.Error
	}
	return storageDirectory.ID, nil
}

func GetStorageDirOwnerByDirId(dirId uint64) (uint, error) {
	var file schema.File
	tx := data.DB.Begin()
	if res := tx.First(&file, "storage_directory_id = ?", dirId); res.Error != nil {
		return 0, res.Error
	}
	return file.UserID, nil
}

func GetChunkDirPath(fileId uint64) (string, error) {
	var directory schema.Directory
	tx := data.DB.Begin()
	res := tx.Table("files").Select(
		"chunk_directories.path").Joins(
		"left join chunk_directories on chunk_directories.id = files.chunk_directory_id").Where(
		"files.id = ?", fileId).Find(
		&directory)
	if res.Error != nil {
		return "", res.Error
	}

	return directory.Path, nil
}

