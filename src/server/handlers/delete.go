package handlers

import (
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

func deleteReassembledFile(user *schema.User, fileID uint) error {
	file, err := repo.GetFileByID(fileID)
	if err != nil {
		return err
	}
	if file.UserID != user.ID {
		return errors.New("access denied")
	}
	return file.Delete(database.DB)
}

func DeleteFile(params files.DeleteFileParams, user *schema.User) middleware.Responder {
	if err := deleteReassembledFile(user, uint(params.FileID)); err != nil {
		return files.NewDeleteFileDefault(500).WithPayload(models.Error("failed to delete the file: " + err.Error()))
	}
	return files.NewDeleteFileOK()
}
