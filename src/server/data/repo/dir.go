package repo

import (
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/scopes"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
)

func GetHomeDir(user *schema.User) (*schema.Directory, error) {
	var homeDir schema.Directory
	tx := data.DB.Scopes(scopes.CurrentUser(user)).Where(&schema.Directory{Name: "home"}).First(&homeDir)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &homeDir, nil
}

func GetSubdirectories(scope scopes.Scope) ([]*models.Directory, error) {
	var directories []*models.Directory
	tx := data.DB.Scopes(scope).
		Model(&schema.Directory{}).
		Not(&schema.Directory{Name: "home"}).
		Find(&directories)
	return directories, tx.Error
}

func GetHomeDirContents(scope scopes.Scope) (*models.DirectoryContents, error) {
	contents, err := GetDirContentsByName(scope, "home")
	if err != nil {
		return contents, err
	}
	contents.Subdirectories, err = GetSubdirectories(scope)
	return contents, err
}

func GetDirectoryByName(scope scopes.Scope, name string) (*schema.Directory, error) {
	var directory schema.Directory
	if err := data.DB.Scopes(scope).Where(&schema.Directory{Name: name}).First(&directory).Error; err != nil {
		return nil, err
	}
	return &directory, nil
}

func GetDirectoryByID(scope scopes.Scope, dirID uint) (*schema.Directory, error) {
	var directory schema.Directory
	if err := data.DB.Scopes(scope).First(&directory, dirID).Error; err != nil {
		return nil, err
	}
	return &directory, nil
}

func GetDirContentsByID(scope scopes.Scope, dirID uint) (*models.DirectoryContents, error) {
	directory, err := GetDirectoryByID(scope, dirID)
	if err != nil {
		return nil, err
	}
	if directory.Name == "home" {
		return GetHomeDirContents(scope)
	}

	var dirContents models.DirectoryContents
	tx := data.DB.Scopes(scope).Table("files").Where(schema.File{StorageDirectoryID: dirID}).Find(&dirContents.Files)
	log.Info(tx.RowsAffected)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &dirContents, nil
}

func GetDirContentsByName(scope scopes.Scope, name string) (*models.DirectoryContents, error) {
	if _, err := GetDirectoryByName(scope, name); err != nil {
		return nil, err
	}
	dir, err := GetDirectoryByName(scope, name)
	if err != nil {
		return nil, err
	}

	var dirContents models.DirectoryContents
	tx := data.DB.Model(&schema.File{}).Where(schema.File{StorageDirectoryID: dir.ID}).Find(&dirContents.Files)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &dirContents, nil
}

func GetDirPathByID (scope scopes.Scope, dirID uint) (string, error) {
	directory, err := GetDirectoryByID(scope, dirID)
	if err != nil {
		return "", err
	}
	return directory.Path, nil
}

func GetChunkDirectoryByID (scope scopes.Scope, dirID uint) (*schema.ChunkDirectory, error) {
	var directory schema.ChunkDirectory
	if err := data.DB.Scopes(scope).First(&directory, dirID).Error; err != nil {
		return nil, err
	}
	return &directory, nil
}

func GetChunkDirPathByFileID (scope scopes.Scope, fileID uint64) (string, error) {
	file, err := GetFileSchemaByID(scope, fileID)

	directory, err := GetChunkDirectoryByID(scope, file.ChunkDirectoryID)
	if err != nil {
		return "", err
	}

	return directory.Path, nil
}

func GetStorageDirectoryPathByFileID (scope scopes.Scope, fileID uint64) (string, error) {
	file, err := GetFileSchemaByID(scope, fileID)
	if err != nil {
		return "", err
	}

	directory, err := GetDirectoryByID(scope, file.StorageDirectoryID)
	return directory.Path, nil
}

func GetStorageDirectoryForFileID(fileId uint) (*schema.Directory, error) {
	var file schema.File
	res := data.DB.Joins("StorageDirectory").Find(&file, fileId)
	if res.Error != nil {
		return nil, res.Error
	}

	return &file.StorageDirectory, nil
}
