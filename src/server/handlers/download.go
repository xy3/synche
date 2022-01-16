package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/xy3/synche/src/files"
	schema2 "github.com/xy3/synche/src/schema"
	"github.com/xy3/synche/src/server"
	"github.com/xy3/synche/src/server/restapi/operations/transfer"
	"path/filepath"
)

// DownloadFile Responds to the client with the file specified in the client's request
func DownloadFile(params transfer.DownloadFileParams, user *schema2.User) middleware.Responder {
	var file schema2.File
	tx := server.DB.Joins("Directory").First(&file, params.FileID)

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

	namedReader := runtime.NamedReader(file.Name, fileReader)

	stat, err := fileReader.Stat()
	if err != nil {
		return transfer.NewDownloadFileNotFound()
	}

	return transfer.NewDownloadFileOK().WithPayload(namedReader).
		WithContentDisposition(fmt.Sprintf("attachment; filename=\"%s\";", file.Name)).
		WithContentLength(uint64(stat.Size()))
}
