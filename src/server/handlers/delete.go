package handlers

import (
	"errors"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

func deleteReassembledFileByID(user *schema.User, fileID uint) error {
	file, err := repo.GetFileByID(fileID, database.DB)
	if err != nil {
		return err
	}

	if file.UserID != user.ID {
		return errors.New("access denied")
	}
	return file.Delete(database.DB)
}

func DeleteFileID(params files.DeleteFileParams, user *schema.User) middleware.Responder {
	if err := deleteReassembledFileByID(user, uint(params.FileID)); err != nil {
		return files.NewDeleteFileDefault(500).WithPayload(models.Error("failed to delete the file: " + err.Error()))
	}
	return files.NewDeleteFileOK()
}

func deleteReassembledFileByPath(path string, user *schema.User) error {
	fullPath, err := repo.BuildFullPath(path, user, database.DB)
	if err != nil {
		return err
	}

	file, err := repo.FindFileByFullPath(fullPath, database.DB)
	if err != nil {
		return err
	}

	if file.UserID != user.ID {
		return errors.New("access denied")
	}
	return file.Delete(database.DB)
}

func DeleteFilePath(params files.DeleteFilepathParams, user *schema.User) middleware.Responder {
	if err := deleteReassembledFileByPath(params.FilePath, user); err != nil {
		return files.NewDeleteFilepathDefault(500).WithPayload(models.Error("failed to delete the file: " + err.Error()))
	}
	log.Infof("deleted file at %v", params.FilePath)
	return files.NewDeleteFilepathOK()
}
