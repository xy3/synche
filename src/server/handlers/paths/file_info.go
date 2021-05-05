package paths

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

func FileInfo(params files.GetFilepathInfoParams, user *schema.User) middleware.Responder {
	file, err := repo.GetFileByPath(user, params.Filepath)
	if err != nil {
		return files.NewGetFilepathInfoNotFound()
	}
	return files.NewGetFilepathInfoOK().WithPayload(repo.ConvertFileToModelsFile(*file))
}
