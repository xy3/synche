package handlers

import (
	"errors"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/xy3/synche/src/server"
	"github.com/xy3/synche/src/server/models"
	"github.com/xy3/synche/src/server/repo"
	"github.com/xy3/synche/src/server/restapi/operations/files"
	"github.com/xy3/synche/src/server/restapi/operations/users"
	"github.com/xy3/synche/src/server/schema"
)

// deleteUserDetails
func deleteUserDetails(email string) error {
	user, err := repo.GetUserByEmail(email, server.DB)
	if err != nil {
		return err
	}

	if err = user.Delete(server.DB); err != nil {
		return err
	}

	return nil
}

// deleteFilesByUserID Deletes all files and directories belonging to a user
func deleteFilesByUserID(userID uint) error {
	directory, err := repo.GetHomeDir(userID, server.DB)
	if err != nil {
		return err
	}

	if err = directory.Delete(true, server.DB); err != nil {
		return err
	}
	return nil
}

// deleteReassembledFileByID Deletes a file specified by its ID
func deleteReassembledFileByID(user *schema.User, fileID uint) error {
	file, err := repo.GetFileByID(fileID, server.DB)
	if err != nil {
		return err
	}

	if file.UserID != user.ID {
		return errors.New("access denied")
	}

	if err := file.Delete(server.DB); err != nil {
		return err
	}

	directory, err := repo.GetDirectoryByID(file.DirectoryID, server.DB)
	if err != nil {
		return err
	}

	if _, err := directory.UpdateFileCount(server.DB); err != nil {
		return err
	}
	return nil
}

// DeleteFileID The handler for deleting files specified by ID. Attempts to delete a file and then responds to the
// client accordingly
func DeleteFileID(params files.DeleteFileParams, user *schema.User) middleware.Responder {
	if err := deleteReassembledFileByID(user, uint(params.FileID)); err != nil {
		return files.NewDeleteFileDefault(500).WithPayload(models.Error("failed to delete the file: " + err.Error()))
	}
	return files.NewDeleteFileOK()
}

// deleteReassembledFileByPath Deletes a file specified by its path
func deleteReassembledFileByPath(path string, user *schema.User) error {
	fullPath, err := repo.BuildFullPath(path, user, server.DB)
	if err != nil {
		return err
	}

	file, err := repo.FindFileByFullPath(fullPath, server.DB)
	if err != nil {
		return err
	}

	if file.UserID != user.ID {
		return errors.New("access denied")
	}

	if err := file.Delete(server.DB); err != nil {
		return err
	}

	directory, err := repo.GetDirectoryByID(file.DirectoryID, server.DB)
	if err != nil {
		return err
	}

	if _, err := directory.UpdateFileCount(server.DB); err != nil {
		return err
	}
	return nil
}

// DeleteFilePath The handler for deleting files specified by their path. Attempts to delete a file and
// then responds to the client accordingly
func DeleteFilePath(params files.DeleteFilepathParams, user *schema.User) middleware.Responder {
	if err := deleteReassembledFileByPath(params.FilePath, user); err != nil {
		return files.NewDeleteFilepathDefault(500).WithPayload(models.Error("failed to delete the file: " + err.Error()))
	}

	log.Infof("deleted file at %v", params.FilePath)
	return files.NewDeleteFilepathOK()
}

// DeleteUser The handler for deleting a user. Deletes all user stored user details and files. It then responds
// to the client accordingly
func DeleteUser(params users.DeleteUserParams, user *schema.User) middleware.Responder {
	if params.Email != user.Email {
		return users.NewDeleteUserDefault(500).WithPayload("failed to authenticate the user")
	}

	if err := deleteFilesByUserID(user.ID); err != nil {
		return users.NewDeleteUserDefault(500).WithPayload(models.Error("failed to delete the user files: " + err.Error()))
	}

	if err := deleteUserDetails(user.Email); err != nil {
		return users.NewDeleteUserDefault(500).WithPayload(models.Error("failed to delete the user details: " + err.Error()))
	}
	return users.NewDeleteUserOK()
}
