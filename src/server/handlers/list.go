package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
)

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

	contents, err := repo.GetDirContentsByID(uint(params.ID))
	if err != nil {
		return errNotFound
	}

	return files.NewListDirectoryOK().WithPayload(contents)
}

func ListHomeDirectory(_ files.ListHomeDirectoryParams, user *schema.User) middleware.Responder {
	homeDirContents, err := repo.GetHomeDirContents(user, database.DB)
	if err != nil {
		return files.NewListHomeDirectoryDefault(500).WithPayload("could not list the home directory")
	}
	return files.NewListHomeDirectoryOK().WithPayload(homeDirContents)
}
