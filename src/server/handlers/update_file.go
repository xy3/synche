package handlers

import (
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"gorm.io/gorm"
	"path/filepath"
)

var (
	ErrDirNotFound = errors.New("directory not found")
)

func updateFile(file *schema.File, user *schema.User, update *models.FileUpdate, db *gorm.DB) (
	newFile *schema.File,
	err error,
) {
	if file.Directory == nil {
		err = database.DB.Preload("Directory").Find(file).Error
		if err != nil || file.Directory == nil {
			return newFile, ErrDirNotFound
		}
	}

	var (
		fullPath  string
		directory = file.Directory
		filename  = file.Name
	)

	// If NewFilePath is set, ignore NewDirectoryID and NewFileName
	if update.NewFilePath != "" {
		fullPath, err = repo.BuildFullPath(update.NewFilePath, user, db)
		err = repo.MoveFile(file, fullPath, db)
		return file, err
	}

	// get directory the file is being moved to by ID
	if update.NewDirectoryID != 0 {
		directory, err = repo.GetDirectoryByID(uint(update.NewDirectoryID), db)
		if err != nil {
			return newFile, ErrDirNotFound
		}
	}

	if update.NewFileName != "" {
		filename = update.NewFileName
	}

	err = repo.MoveFile(file, filepath.Join(directory.Path, filename), db)
	return file, err
}

func UpdateFileByID(params files.UpdateFileByIDParams, user *schema.User) middleware.Responder {
	file, err := repo.GetFileByID(uint(params.FileID), database.DB)
	if err != nil {
		return files.NewUpdateFileByIDDefault(404).WithPayload("file not found")
	}

	if file.UserID != user.ID {
		return files.NewUpdateFileByIDUnauthorized()
	}

	newFile, err := updateFile(file, user, params.FileUpdate, database.DB)
	if err != nil {
		return files.NewUpdateFileByIDDefault(500).WithPayload(models.Error("failed to update the file: " + err.Error()))
	}

	return files.NewUpdateFileByIDOK().WithPayload(newFile.ConvertToFileModel())
}

func UpdateFileByPath(params files.UpdateFileByPathParams, user *schema.User) middleware.Responder {
	var (
		file    *schema.File
		newFile *schema.File
		err404  = files.NewUpdateFileByPathDefault(404)
		err500  = files.NewUpdateFileByPathDefault(500)
	)

	fullPath, err := repo.BuildFullPath(params.FilePath, user, database.DB)
	if err != nil {
		return err500.WithPayload(models.Error(err.Error()))
	}

	file, err = repo.FindFileByFullPath(fullPath, database.DB)
	if err != nil {
		return err404.WithPayload("file not found")
	}

	if file.UserID != user.ID {
		return files.NewUpdateFileByPathUnauthorized()
	}

	newFile, err = updateFile(file, user, params.FileUpdate, database.DB)
	if err != nil {
		return err500.WithPayload(models.Error("failed to update the file: " + err.Error()))
	}

	return files.NewUpdateFileByPathOK().WithPayload(newFile.ConvertToFileModel())
}
