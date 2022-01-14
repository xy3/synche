package repo

import (
	"github.com/xy3/synche/src/hash"
	"github.com/xy3/synche/src/server/models"
	"github.com/xy3/synche/src/server/schema"
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

	return GetDirContentsByID(homeDir.ID, db)
}

func GetDirContentsByID(dirID uint, db *gorm.DB) (*models.DirectoryContents, error) {
	contents := &models.DirectoryContents{}
	directory := &schema.Directory{}

	tx := db.Where("id = ?", dirID).First(directory)
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

	tx = db.Where(&schema.Directory{ParentID: &dirID}).Find(&contents.SubDirectories)
	if tx.Error != nil {
		return nil, tx.Error
	}

	tx = db.Where(&schema.File{DirectoryID: dirID}).Find(&contents.Files)
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
