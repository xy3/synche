package repo

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/models"
	"gorm.io/gorm"
)

func GetHomeDirContents(user *schema.User, db *gorm.DB) (*models.DirectoryContents, error) {
	var (
		err     error
		homeDir *schema.Directory
	)
	homeDir, err = GetHomeDir(user.ID, db)

	if err != nil {
		if err.Error() == "record not found" {
			homeDir, err = SetupUserHomeDir(user, db)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return GetDirContentsByID(homeDir.ID)
}

func GetDirContentsByID(dirID uint) (*models.DirectoryContents, error) {
	contents := &models.DirectoryContents{}

	directory := &schema.Directory{}

	tx := database.DB.Where("id = ?", dirID).First(directory)
	if tx.Error != nil {
		return nil, tx.Error
	}

	contents.CurrentDir = &models.Directory{
		FileCount: uint64(directory.FileCount),
		ID:        uint64(directory.ID),
		Name:      directory.Name,
		Path:      directory.Path,
		PathHash:  directory.PathHash,
	}

	if directory.ParentID != nil {
		contents.CurrentDir.ParentDirectoryID = uint64(*directory.ParentID)
	}

	tx = database.DB.Where(&schema.Directory{ParentID: &dirID}).Find(&contents.SubDirectories)
	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = database.DB.Where(&schema.File{DirectoryID: dirID}).Find(&contents.Files)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return contents, nil
}

func GetDirWithContentsFromPath(path string, db *gorm.DB) (*schema.Directory, error) {
	pathHash := hash.PathHash(path)
	directory := &schema.Directory{}
	tx := db.Preload("Children").Preload("Files").Where(schema.Directory{PathHash: pathHash}).First(directory)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return directory, nil
}
