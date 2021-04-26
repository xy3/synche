package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	f "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/scopes"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"path/filepath"
)

func updateFileStorageDirectory(user *schema.User, file *models.File, dirID uint) (*models.File, error) {
	file, err := repo.UpdateFileStorageDirectory(scopes.CurrentUser(user), *file.ID, dirID)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func getPaths(user *schema.User, file *models.File, directory *schema.Directory) (oldPath string, newPath string, err error) {
	// get current path to file
	currentDirPath, err := repo.GetDirPathByID(scopes.CurrentUser(user), uint(*file.StorageDirectoryID))
	if err != nil {
		return "", "", err
	}

	oldPath = filepath.Join(currentDirPath, *file.Name)

	// new file path
	newDirPath, err := repo.GetDirPathByID(scopes.CurrentUser(user), directory.ID)
	if err != nil {
		return "", "", err
	}

	newPath = filepath.Join(newDirPath, *file.Name)

	return oldPath, newPath, nil
}

func moveFileJob(user *schema.User, file *models.File, directory *schema.Directory) error {
	oldPath, newPath, err := getPaths(user, file, directory)
	if err != nil {
		return err
	}

	if err := f.AppFS.Rename(oldPath, newPath); err != nil {
		log.Infof("old path: %v, new path: %v, error %v", oldPath, newPath, err)
		return err
	}

	return nil
}

func getDirID (user *schema.User, params files.MoveFileParams) (uint, error) {
	if params.DirectoryName != nil {
		directory, err := repo.GetDirectoryByName(scopes.CurrentUser(user), *params.DirectoryName)
		if err != nil {
			return 0, err
		}
		return directory.ID, nil
	}
	return uint(*params.DirectoryID), nil

}

func MoveFile(
	params files.MoveFileParams,
	user *schema.User,
	) middleware.Responder {

	dirID, err := getDirID(user, params)
	if err != nil {
		return files.NewMoveFileDefault(404).WithPayload("invalid directory specification")
	}

	file, err := repo.GetFileModelByID(scopes.CurrentUser(user), params.FileID)
	if err != nil {
		return files.NewMoveFileDefault(404).WithPayload("file not found")
	}

	directory, err := repo.GetDirectoryByID(scopes.CurrentUser(user), dirID)
	if err != nil {
		return files.NewMoveFileDefault(404).WithPayload("directory not found")
	}

	if directory.ID == dirID {
		return files.NewMoveFileDefault(404).WithPayload("the file is already in this directory")
	}

	// move file
	if err = moveFileJob(user, file, directory); err != nil {
		return files.NewMoveFileDefault(404).WithPayload("error trying to move file")
	}

	// update db
	updatedFile, err := updateFileStorageDirectory(user, file, directory.ID)
	if err != nil {
		log.WithError(err).Error("File information could not be updated")
	}

	return files.NewMoveFileOK().WithPayload(updatedFile)
}
