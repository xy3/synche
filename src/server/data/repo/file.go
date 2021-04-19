package repo

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
)

func GetFilenameByFileId(fileId uint64) (string, error) {
	var file schema.File
	tx := data.DB.Begin()
	if res := tx.First(&file, fileId); res.Error != nil {
		return "", res.Error
	}
	return file.Name, nil
}

func GetFileOwnerByFileId(fileId uint64) (uint, error) {
	var file schema.File
	tx := data.DB.Begin()
	if res := tx.First(&file, fileId); res.Error != nil {
		return 0, res.Error
	}
	return file.UserID, nil
}

func UpdateFileName(uploadRequestId uint64, uniqueFilename string) {
	var file schema.File
	tx := data.DB.Begin()
	tx.Model(&file).Where("id", uploadRequestId).Update("name", uniqueFilename)
	tx.Commit()
}
