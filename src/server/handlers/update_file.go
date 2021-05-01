package handlers

import (
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	f "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
)

func ConvertToFileModel(file *schema.File) *models.File {
	return &models.File{
		ID:          uint64(file.ID),
		Name:        file.Name,
		Size:        file.Size,
		Hash:        file.Hash,
		DirectoryID: uint64(file.DirectoryID),
		Available:   file.Available,
	}
}

func UpdateFile(params files.UpdateFileParams, user *schema.User) middleware.Responder {
	var (
		err                      error
		newFile                  *schema.File
		errNotFound              = files.NewUpdateFileDefault(404).WithPayload("file not found")
		errFileDirNotFound       = files.NewUpdateFileDefault(404).WithPayload("file directory not found")
		errRenameFailed          = files.NewUpdateFileDefault(500).WithPayload("failed to rename the file")
		errUpdateDirectoryFailed = files.NewUpdateFileDefault(500).WithPayload("failed to update the file's directory")
	)

	file := &schema.File{}
	tx := database.DB.Where(schema.File{
		Model:  gorm.Model{ID: uint(params.FileID)},
		UserID: user.ID,
	}).First(file)

	if tx.Error != nil {
		return errNotFound
	}

	if file.Directory == nil {
		if err = database.DB.Preload("Directory").Find(file).Error; err != nil {
			return errFileDirNotFound
		}
		if file.Directory == nil {
			return errFileDirNotFound
		}
	}

	if params.Filename != nil && file.Name != *params.Filename {
		newFile, err = updateFilename(file, *params.Filename)
		if err != nil {
			return errRenameFailed
		}
	}

	if params.DirectoryID != nil && file.DirectoryID != uint(*params.DirectoryID) {
		newFile, err = updateDirectory(file, uint(*params.DirectoryID))
		if err != nil {
			return errUpdateDirectoryFailed
		}
	}

	return files.NewUpdateFileOK().WithPayload(ConvertToFileModel(newFile))
}

func moveFileToNewDir(file *schema.File, newDirID uint) error {
	var directory schema.Directory
	tx := database.DB.Where("id = ?", newDirID).Find(&directory)
	if tx.Error != nil {
		return tx.Error
	}
	filePath, err := file.Path(database.DB)
	if err != nil {
		return err
	}
	newPath := filepath.Join(directory.Path, file.Name)
	if err = f.Afs.Rename(filePath, newPath); err != nil {
		return err
	}
	return nil
}

func updateDirectory(file *schema.File, newDirID uint) (*schema.File, error) {
	if err := moveFileToNewDir(file, newDirID); err != nil {
		return nil, err
	}

	file.DirectoryID = newDirID
	tx := database.DB.Save(&file)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return file, nil
}

func updateFilename(file *schema.File, newName string) (*schema.File, error) {
	if strings.ContainsAny(newName, "/\\") {
		return file, errors.New("file names cannot contain slashes")
	}

	filePath, err := file.Path(database.DB)
	if err != nil {
		return nil, err
	}

	newPath := filepath.Join(file.Directory.Path, file.Name)
	if err = f.Afs.Rename(filePath, newPath); err != nil {
		return nil, err
	}

	file.Name = newName
	tx := database.DB.Save(&file)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return file, nil
}
