package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/scopes"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/jobs"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

func convertToModelsDir(directory *schema.Directory) *models.Directory {
	return &models.Directory{
		Name: directory.Name,
		ID:   uint64(directory.ID),
	}
}

func CreateDirectory(params files.CreateDirectoryParams, user *schema.User) middleware.Responder {
	if len(params.DirectoryName) < 3 {
		return files.NewCreateDirectoryDefault(400).WithPayload("directory names must be greater than 3 characters in length")
	}

	newDir, err := jobs.CreateDirectoryInsideUserHomeDir(params.DirectoryName, user)
	if err != nil {
		return files.NewCreateDirectoryDefault(500).WithPayload("could not create the directory")
	}

	var directory *schema.Directory
	directory, err = repo.GetDirectoryByName(scopes.CurrentUser(user), params.DirectoryName)
	if err == nil {
		return files.NewCreateDirectoryOK().WithPayload(convertToModelsDir(directory))
	}

	db := data.DB.Begin()

	var userDir = schema.Directory{
		Name:     params.DirectoryName,
		Path:     newDir,
		PathHash: hash.MD5HashString(newDir),
		UserID:   user.ID,
	}

	if err = db.Create(&userDir).Error; err != nil {
		db.Rollback()
		return files.NewCreateDirectoryDefault(500).WithPayload("could not store the directory in the database")
	}

	db.Commit()

	return files.NewCreateDirectoryOK().WithPayload(convertToModelsDir(&userDir))
}

func DeleteDirectory(params files.DeleteDirectoryParams, user *schema.User) middleware.Responder {
	log.Info("Deleting dir: ", params.ID)
	// This handler does not ask for confirmation. The directory is completely gone if this handler is called.
	db := data.DB.Begin()
	// This scope just makes sure "where user_id = the users ID" is applied
	scope := scopes.CurrentUser(user)
	db = db.Scopes(scope)

	directory, err := repo.GetDirectoryByID(scope, uint(params.ID))
	if err != nil {
		return files.NewDeleteDirectoryDefault(404).WithPayload("directory not found")
	}

	if directory.Name == "home" {
		return files.NewDeleteDirectoryDefault(400).WithPayload("you cannot delete your home directory")
	}

	if err = db.Where(schema.File{StorageDirectoryID: uint(params.ID)}).Delete(&schema.File{}).Error; err != nil {
		db.Rollback()
		return files.NewDeleteDirectoryDefault(500).WithPayload(models.Error("failed to delete the directory: " + err.Error()))
	}

	if err = db.Delete(&schema.Directory{}, params.ID).Error; err != nil {
		db.Rollback()
		return files.NewDeleteDirectoryDefault(500).WithPayload(models.Error("failed to delete the directory: " + err.Error()))
	}

	if err = jobs.DeleteDirectory(directory.Path); err != nil {
		db.Rollback()
		return files.NewDeleteDirectoryDefault(500).WithPayload(models.Error("failed to delete the directory: " + err.Error()))
	}
	db.Commit()

	return files.NewDeleteDirectoryOK()
}

func DirectoryInfo(params files.GetDirectoryInfoParams, user *schema.User) middleware.Responder {
	directory, err := repo.GetDirectoryByID(scopes.CurrentUser(user), uint(params.ID))
	if err != nil {
		return files.NewGetDirectoryInfoDefault(404).WithPayload("directory not found")
	}
	return files.NewGetDirectoryInfoOK().WithPayload(convertToModelsDir(directory))
}
