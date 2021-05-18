package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

// findExistingDirByParentID Returns a directory specified by it's parent directory ID
func findExistingDirByParentID(dirName string, parentDirID uint) (*schema.Directory, error) {
	var directory schema.Directory
	tx := database.DB.Where(schema.Directory{Name: dirName, ParentID: &parentDirID}).First(&directory)
	return &directory, tx.Error
}

// CreateDirectory Creates a directory on disk and updates the database
func CreateDirectory(params files.CreateDirectoryParams, user *schema.User) middleware.Responder {
	var (
		err             error
		directory       *schema.Directory
		homeDir         *schema.Directory
		defaultRes      = files.NewCreateDirectoryDefault
		errDirTooShort  = defaultRes(400).WithPayload("directory names must be greater than 3 characters in length")
		errCreateFailed = defaultRes(500).WithPayload("could not create the directory")
		errNoParentDir  = defaultRes(500).WithPayload("could not locate parent directory")
	)

	if len(params.DirectoryName) < 3 {
		return errDirTooShort
	}

	var parentDirID uint
	if params.ParentDirectoryID != nil {
		parentDirID = uint(*params.ParentDirectoryID)
	} else {
		homeDir, err = repo.GetHomeDir(user.ID, database.DB)
		if err != nil {
			return errNoParentDir
		}
		parentDirID = homeDir.ID
	}

	directory, err = findExistingDirByParentID(params.DirectoryName, parentDirID)

	if err != nil {
		directory, err = repo.CreateDirectory(params.DirectoryName, parentDirID, database.DB)
		if err != nil {
			return errCreateFailed
		}
	}

	modelsDir := directory.ConvertToModelsDir()
	return files.NewCreateDirectoryOK().WithPayload(modelsDir)
}

// DeleteDirectory Deletes a directory from disk and from the database
func DeleteDirectory(params files.DeleteDirectoryParams, user *schema.User) middleware.Responder {
	var (
		err         error
		directory   *schema.Directory
		errNotFound = files.NewDeleteDirectoryDefault(404).WithPayload("directory not found")
		errNoAccess = files.NewDeleteDirectoryDefault(501).WithPayload("you do not have access to this directory")
		err500      = files.NewDeleteDirectoryDefault(500)
	)

	log.Info("Deleting dir: ", params.ID)
	// This handler does not ask for confirmation. The directory is completely gone if this handler is called.
	// This scope just makes sure "where user_id = the users ID" is applied
	directory, err = repo.GetDirectoryByID(uint(params.ID), database.DB)
	if err != nil {
		return errNotFound
	}

	if directory.UserID != user.ID {
		return errNoAccess
	}

	if err = directory.Delete(true, database.DB); err != nil {
		return err500.WithPayload(models.Error("failed to delete the directory: " + err.Error()))
	}

	return files.NewDeleteDirectoryOK()
}
