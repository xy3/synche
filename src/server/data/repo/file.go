package repo

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/scopes"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
)

func GetFileModelByID(scope scopes.Scope, fileID uint64) (*models.File, error) {
	var file models.File
	res := data.DB.Model(&schema.File{}).Scopes(scope).Find(&file, fileID)
	if res.Error != nil {
		return nil, res.Error
	}
	return &file, nil
}

func GetFileSchemaByID (scope scopes.Scope, fileID uint64) (*schema.File, error) {
	var file schema.File
	if err := data.DB.Scopes(scope).First(&file, fileID).Error; err != nil {
	return nil, err
	}
	return &file, nil
}

func GetFilenameByFileID(scope scopes.Scope, fileID uint64) (string, error) {
	file, err := GetFileModelByID(scope, fileID)
	if err != nil {
		return "", err
	}
	return *file.Name, nil
}

func UpdateFileDeletedAt(fileID uint64) error {
	tx := data.DB.Begin()
	tx.Delete(&schema.File{}, fileID)
	tx.Commit()
	if tx.Error != nil {
		return tx.Error
	}
	return nil
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

func UpdateFileStorageDirectory(scope scopes.Scope, fileID uint64, dirID uint) (*models.File, error) {
	var file models.File
	tx := data.DB.Begin()
	tx.Model(&schema.File{}).Where("id", fileID).Update("StorageDirectoryID", dirID)
	tx.Commit()
	tx = data.DB.Model(&schema.File{}).Scopes(scope).Find(&file, fileID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &file, nil
}