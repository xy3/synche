package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/scopes"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"strings"
)

func getDirContents(params files.ListDirectoryParams, user *schema.User) (*models.DirectoryContents, error) {
	if params.DirectoryID != nil {
		return repo.GetDirContentsByID(scopes.CurrentUser(user), uint(*params.DirectoryID))
	}
	return repo.GetDirContentsByName(scopes.CurrentUser(user), *params.DirectoryName)
}

func listHome(user *schema.User) middleware.Responder {
	homeDirContents, err := repo.GetHomeDirContents(scopes.CurrentUser(user))
	if err != nil {
		return files.NewListDirectoryDefault(500).WithPayload("could not list the home directory")
	}
	return files.NewListDirectoryOK().WithPayload(homeDirContents)
}

func ListDirectory(params files.ListDirectoryParams, user *schema.User) middleware.Responder {
	// List the home directory if either no parameters were provided, or if the name provided is "home"
	if params.DirectoryID == nil && params.DirectoryName == nil {
		return listHome(user)
	}
	if params.DirectoryName != nil && strings.ToLower(*params.DirectoryName) == "home" {
		return listHome(user)
	}

	contents, err := getDirContents(params, user)
	if err != nil {
		return files.NewDeleteDirectoryDefault(404).WithPayload("directory not found")
	}

	return files.NewListDirectoryOK().WithPayload(contents)
}
