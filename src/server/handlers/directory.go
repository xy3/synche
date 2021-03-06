package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"path/filepath"
	"strings"
)

// findExistingDirByParentID Returns a directory specified by it's parent directory ID
func findExistingDirByParentID(dirName string, parentDirID uint) (*schema.Directory, error) {
	var directory schema.Directory
	tx := database.DB.Where(schema.Directory{Name: dirName, ParentID: &parentDirID}).First(&directory)
	return &directory, tx.Error
}

// CreateDirectoryByPath Creates a directory on disk and updates the database, location is specified
// by the path to the directory
func CreateDirectoryByPath(params files.CreateDirPathParams, user *schema.User) middleware.Responder {
	var (
		err             error
		parentPath      string
		parentDir       *schema.Directory
		directory       *schema.Directory
		defaultRes      = files.NewCreateDirPathDefault
		errDirTooShort  = defaultRes(400).WithPayload("directory names must be greater than 3 characters in length")
		errCreateFailed = defaultRes(500).WithPayload("could not create the directory")
		errNoParentDir  = defaultRes(500).WithPayload("could not locate parent directory")
	)

	if len(filepath.Base(params.DirPath)) < 3 {
		return errDirTooShort
	}

	// if parentPath == "." then create in the home directory
	trimmedPath := filepath.Dir(strings.TrimRight(strings.TrimSpace(params.DirPath), "/"))
	if trimmedPath == "." {
		parentDir, err = repo.GetHomeDir(user.ID, database.DB)
	} else {
		parentPath, err = repo.BuildFullPath(trimmedPath, user, database.DB)
		if err != nil {
			return errNoParentDir
		}
		parentDir, err = repo.GetDirByPath(parentPath, database.DB)
		if err != nil {
			return errNoParentDir
		}
	}

	directory, err = findExistingDirByParentID(filepath.Base(params.DirPath), parentDir.ID)

	if err != nil {
		directory, err = repo.CreateDirectory(filepath.Base(params.DirPath), parentDir.ID, database.DB)
		if err != nil {
			return errCreateFailed
		}
	}

	modelsDir := directory.ConvertToModelsDir()
	return files.NewCreateDirPathOK().WithPayload(modelsDir)
}

// CreateDirectory Creates a directory on disk and updates the database, location is specified
// by the parent directory's ID
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

// DeleteDirectoryByPath Deletes a directory on disk and updates the database, location is specified
// by the path to the directory
func DeleteDirectoryByPath(params files.DeleteDirPathParams, user *schema.User) middleware.Responder {
	var (
		err         error
		path        string
		directory   *schema.Directory
		errNotFound = files.NewDeleteDirPathDefault(404).WithPayload("directory not found")
		err500      = files.NewDeleteDirPathDefault(500)
		errNoAccess = files.NewDeleteDirPathUnauthorized()
	)

	log.Info("Deleting dir: ", params.DirPath)
	trimmedPath := strings.TrimRight(strings.TrimSpace(params.DirPath), "/")
	path, err = repo.BuildFullPath(trimmedPath, user, database.DB)

	directory, err = repo.GetDirByPath(path, database.DB)
	if err != nil {
		return errNotFound
	}

	if directory.UserID != user.ID {
		return errNoAccess
	}

	if err = directory.Delete(true, database.DB); err != nil {
		return err500.WithPayload(models.Error("failed to delete the directory: " + err.Error()))
	}

	return files.NewDeleteDirPathOK()
}

// DeleteDirectory Deletes a directory from disk and from the database, location is specified
// by the parent directory's ID
func DeleteDirectory(params files.DeleteDirectoryParams, user *schema.User) middleware.Responder {
	var (
		err         error
		directory   *schema.Directory
		errNotFound = files.NewDeleteDirectoryDefault(404).WithPayload("directory not found")
		err500      = files.NewDeleteDirectoryDefault(500)
		errNoAccess = files.NewDeleteDirectoryUnauthorized()
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

// ListDirectoryByPath Retrieves the contents of the specified directory and replies to the client.
// Directory is specified by its path
func ListDirectoryByPath(params files.ListDirPathInfoParams, user *schema.User) middleware.Responder {
	var (
		err         error
		path        string
		directory   *schema.Directory
		errNotFound = files.NewListDirPathInfoDefault(404).WithPayload("directory not found")
		errNoAccess = files.NewListDirPathInfoUnauthorized()
	)

	trimmedPath := strings.TrimRight(strings.TrimSpace(params.DirPath), "/")
	path, err = repo.BuildFullPath(trimmedPath, user, database.DB)

	directory, err = repo.GetDirByPath(path, database.DB)
	if err != nil {
		return errNotFound
	}

	if directory.UserID != user.ID {
		return errNoAccess
	}

	contents, err := repo.GetDirContentsByID(directory.ID, database.DB)
	if err != nil {
		return errNotFound
	}

	return files.NewListDirPathInfoOK().WithPayload(contents)
}

// ListDirectory Retrieves the contents of the specified directory and replies to the client.
// Directory is specified by ID
func ListDirectory(params files.ListDirectoryParams, user *schema.User) middleware.Responder {
	var (
		err         error
		directory   *schema.Directory
		errNotFound = files.NewListDirectoryDefault(404).WithPayload("directory not found")
	)

	directory, err = repo.GetDirectoryByID(uint(params.ID), database.DB)
	if err != nil {
		return errNotFound
	}

	if directory.UserID != user.ID {
		return files.NewListDirectoryUnauthorized()
	}

	contents, err := repo.GetDirContentsByID(uint(params.ID), database.DB)
	if err != nil {
		return errNotFound
	}

	return files.NewListDirectoryOK().WithPayload(contents)
}

// ListHomeDirectory Retrieves the contents of the home directory and replies with these to the client
func ListHomeDirectory(_ files.ListHomeDirectoryParams, user *schema.User) middleware.Responder {
	homeDirContents, err := repo.GetHomeDirContents(user, database.DB)
	if err != nil {
		return files.NewListHomeDirectoryDefault(500).WithPayload("could not list the home directory")
	}
	return files.NewListHomeDirectoryOK().WithPayload(homeDirContents)
}
