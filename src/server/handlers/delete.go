package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	f "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"path/filepath"
)

func deleteReassembledFile(fileId uint64) error {
	filePath := c.Config.Server.StorageDir
	filename, err := repo.GetFilenameByFileId(fileId)
	if err != nil {
		return err
	}

	reassembledFileLocation := filepath.Join(filePath, filename)
	if err := f.AppFS.Remove(reassembledFileLocation); err != nil {
		return err
	}
	return nil
}

func deleteChunkDir(fileId uint64) error {
	dirPath, err := repo.GetChunkDirPath(fileId)
	if err != nil {
		return err
	}
	if err := f.Afs.RemoveAll(dirPath); err != nil {
		return err
	}
	return nil
}

func isFileOwner(fileId uint64, userId uint) (bool, error) {
	fileOwner, err := repo.GetFileOwnerByFileId(fileId)
	if fileOwner == userId && err == nil {
		return true, nil
	}
	return false, err
}

func DeleteFile(
	params files.DeleteFileParams,
	user *schema.User,
) middleware.Responder {
	fileId := params.FileID
	// check if the user owns the file. this can adapted to check if the user is an admin either
	if userPermission, err := isFileOwner(fileId, user.ID); userPermission == true && err == nil {
		if err := deleteChunkDir(fileId); err != nil {
			return files.NewDeleteFileNotFound()
		}
		if err := deleteReassembledFile(fileId); err != nil {
			return files.NewDeleteFileNotFound()
		}
	} else {
		return files.NewDeleteFileUnauthorized()
	}
	log.Infof("File deleted with id: %v", fileId)
	return files.NewDeleteFileCreated().WithPayload(&models.FileDeleted{FileID: fileId})
}
