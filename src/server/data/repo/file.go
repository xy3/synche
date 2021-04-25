package repo

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
)

func GetFileByID(fileID uint64) (*schema.File, error) {
	var file schema.File
	res := data.DB.First(&file, fileID)
	if res.Error != nil {
		return nil, res.Error
	}
	return &file, nil
}

func GetFilenameByFileId(fileID uint64) (string, error) {
	file, err := GetFileByID(fileID)
	if err != nil {
		return "", err
	}
	return file.Name, nil
}

func GetFileOwnerByFileId(fileID uint64) (uint, error) {
	file, err := GetFileByID(fileID)
	if err != nil {
		return 0, err
	}
	return file.UserID, nil
}

func UpdateFilenameForUploadID(uploadID uint, newFilename string) error {
	var upload schema.Upload
	tx := data.DB.Begin()
	if err := tx.First(&upload, uploadID).Error; err != nil {
		return err
	}
	var file schema.File
	tx.Model(&file).Where("id", upload.FileID).Update("name", newFilename)
	tx.Commit()
	return nil
}
