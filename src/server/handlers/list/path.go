package list

import (
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"strings"
)

func ByDirPath(
	params files.ListDPathParams,
	user *schema.User,
) middleware.Responder {
	dirPath := params.DirPath
	if strings.ToLower(dirPath) == "home" {
		dirPath = c.Config.Server.StorageDir
	}

	dirID, err := repo.GetStorageDirIdByPath(dirPath)
	if err != nil {
		return files.NewListDPathUnauthorized()
	}

	if userPermission, err := isDirOwner(uint64(dirID), user.ID); userPermission == true && err == nil {
		contents, err := repo.GetStorageDirContents(uint64(dirID))
		if err != nil {
			return files.NewListDPathUnauthorized()
		}

		dirContents, err := convertToFileModel(contents)
		if err != nil {
			return files.NewListDPathUnauthorized()
		}

		log.Infof("Returning contents for directory with ID: %v", dirPath)
		return files.NewListDPathOK().WithPayload(&models.DirectoryContents{Contents: dirContents})
	}
	return files.NewListDPathUnauthorized()
}
