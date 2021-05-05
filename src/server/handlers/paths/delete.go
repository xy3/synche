package paths

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

func deleteReassembledFileByPath(user *schema.User, path string) error {
	file, err := repo.GetFileByPath(user, path)
	if err != nil {
		return err
	}

	if file.UserID != user.ID {
		return errors.New("access denied")
	}
	return file.Delete(database.DB)
}

func DeleteFilepath(params files.DeleteFilepathParams, user *schema.User) middleware.Responder {
	if err := deleteReassembledFileByPath(user, params.Filepath); err != nil {
		return files.NewDeleteFilepathDefault(500).WithPayload(models.Error("failed to delete the file: " + err.Error()))
	}
	log.Infof("deleted file at %v", params.Filepath)
	return files.NewDeleteFilepathOK()
}
