package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	f "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/scopes"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations/files"
	"path/filepath"
)

func deleteReassembledFile(user *schema.User, fileID uint64) error {
	filename, err := repo.GetFilenameByFileID(scopes.CurrentUser(user), fileID)
	if err != nil {
		return err
	}

	dirPath, err := repo.GetStorageDirectoryPathByFileID(scopes.CurrentUser(user), fileID)
	if err != nil {
		return err
	}

	reassembledFileLocation := filepath.Join(dirPath, filename)
	if err := f.AppFS.Remove(reassembledFileLocation); err != nil {
		return err
	}
	return nil
}

func deleteChunkDir(user *schema.User, fileId uint64) error {
	dirPath, err := repo. GetChunkDirPathByFileID(scopes.CurrentUser(user), fileId)
	if err != nil {
		return err
	}
	if err := f.Afs.RemoveAll(dirPath); err != nil {
		return err
	}
	return nil
}

func DeleteFile(
	params files.DeleteFileParams,
	user *schema.User,
	) middleware.Responder {
	fileID := params.FileID
	if err := deleteChunkDir(user, fileID); err != nil {
		return files.NewDeleteFileNotFound()
	}

	if err := deleteReassembledFile(user, fileID); err != nil {
		return files.NewDeleteFileNotFound()
	}

	if err := repo.UpdateFileDeletedAt(fileID); err != nil {
		return files.NewDeleteFileUnauthorized()
	}

	return files.NewDeleteFileCreated().WithPayload(&models.FileDeleted{FileID: fileID})
}
