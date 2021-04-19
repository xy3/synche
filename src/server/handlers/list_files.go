package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

func convertToFileModel(directoryContents []schema.File) ([]*models.File, error) {
	var contents []*models.File
	err := copier.Copy(&contents, &directoryContents)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func isDirOwner(dirId uint64, userId uint) (bool, error) {
	dirOwner, err := repo.GetDirOwnerByDirId(dirId)
	if dirOwner == userId && err == nil {
		return true, nil
	}
	return false, err
}

func ListFiles(
	params files.ListParams,
	user *schema.User,
) middleware.Responder {
	dirId := params.DirectoryID
	if userPermission, err := isDirOwner(dirId, user.ID); userPermission == true && err == nil {
		contents, err := repo.GetDirContents(dirId)
		if err != nil {
			return files.NewListUnauthorized()
		}

		dirContents, err := convertToFileModel(contents)
		if err != nil {
			return files.NewListUnauthorized()
		}

		log.Infof("Returning contents for directory with ID: %v", dirId)
		return files.NewListOK().WithPayload(&models.DirectoryContents{Contents: dirContents})
	}
	return files.NewListUnauthorized()
}
