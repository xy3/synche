package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/transfer"
	"path/filepath"
)

func DownloadFile(params transfer.DownloadFileParams, user *schema.User) middleware.Responder {
	var file schema.File
	tx := database.DB.Joins("Directory").First(&file, params.FileID)

	if tx.Error != nil {
		return transfer.NewDownloadFileNotFound()
	}

	if file.UserID != user.ID {
		return transfer.NewDownloadFileForbidden()
	}

	filePath := filepath.Join(file.Directory.Path, file.Name)
	fileReader, err := files.Afs.Open(filePath)
	log.Debugf("Reading file: %s", filePath)
	if err != nil {
		return transfer.NewDownloadFileDefault(500).WithPayload("failed to read the file")
	}
	stat, err := fileReader.Stat()

	if err != nil {
		return transfer.NewDownloadFileNotFound()
	}

	return transfer.NewDownloadFileOK().WithPayload(fileReader).WithContentLength(uint64(stat.Size()))
}
