package repo

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
)

func GetDirContents(dirId uint64) ([]schema.File, error) {
	var files []schema.File
	tx := data.DB.Find(&files, schema.File{DirectoryID: uint(dirId)})
	if tx.Error != nil {
		return files, tx.Error
	}
	return files, nil
}

func GetDirPath(fileId uint64) (dirPath string, err error) {
	var directory schema.Directory
	tx := data.DB.Begin()
	res := tx.Table("files").Select(
		"directories.path").Joins(
		"left join directories on directories.id = files.directory_id").Where(
		"files.id = ?", fileId).Find(
		&directory)
	if res.Error != nil {
		return "", res.Error
	}

	return  directory.Path, nil
}

func GetDirOwnerByDirId(dirId uint64) (uint, error) {
	var file schema.File
	tx := data.DB.Begin()
	if res := tx.First(&file, dirId); res.Error != nil {
		return 0, res.Error
	}
	return file.UserID, nil
}


