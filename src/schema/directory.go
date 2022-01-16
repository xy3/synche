package schema

import (
	"errors"
	"github.com/go-openapi/strfmt"
	"github.com/xy3/synche/src/files"
	"github.com/xy3/synche/src/models"
	"gorm.io/gorm"
)

var (
	ErrDirNotEmpty = errors.New("directory is not empty")
)

// swagger:model Directory
type Directory struct {
	Model
	Name      string `gorm:"not null"`
	Path      string `gorm:"not null"`
	PathHash  string `gorm:"size:32;uniqueIndex"`
	FileCount int64
	UserID    uint
	ParentID  *uint

	Parent   *Directory  `gorm:"foreignKey:id;association_foreignKey:parent_id;association_autoupdate:false;association_autocreate:false"`
	User     User        `gorm:"association_autoupdate:false;association_autocreate:false"`
	Children []Directory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:parent_id;association_autoupdate:false;association_autocreate:false"`
	Files    []File      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:directory_id;association_autoupdate:false;association_autocreate:false"`
}

// Validate validates this directory
func (m *Directory) Validate(formats strfmt.Registry) error {
	return nil
}

// executeDelete removes the directory and everything contained within it from the disk
func (dir *Directory) executeDelete() error {
	return files.Afs.RemoveAll(dir.Path)
}

func (dir *Directory) Delete(forceDelete bool, db *gorm.DB) (err error) {
	if err = dir.executeDelete(); err != nil {
		return err
	}

	var filesInDir int64
	db.Model(File{}).Where(File{DirectoryID: dir.ID}).Count(&filesInDir)
	if filesInDir != 0 && !forceDelete {
		return ErrDirNotEmpty
	}

	if err = db.Where(File{DirectoryID: dir.ID}).Delete(&File{}).Error; err != nil {
		return err
	}

	if err = db.Unscoped().Delete(&Directory{}, dir).Error; err != nil {
		return err
	}

	return nil
}

func (dir *Directory) UpdateFileCount(db *gorm.DB) (num int64, err error) {
	tx := db.Model(&File{}).Where("directory_id = ?", dir.ID).Count(&num)
	if tx.Error != nil {
		return num, tx.Error
	}
	dir.FileCount = num
	return num, db.Save(dir).Error
}

// ConvertToModelsDir Translates a schema directory to a model directory
func (dir *Directory) ConvertToModelsDir() *models.Directory {
	var parentDirID uint64
	if dir.ParentID != nil {
		parentDirID = uint64(*dir.ParentID)
	}
	return &models.Directory{
		FileCount:         uint64(dir.FileCount),
		ID:                uint64(dir.ID),
		Name:              dir.Name,
		ParentDirectoryID: parentDirID,
		Path:              dir.Path,
		PathHash:          dir.PathHash,
	}
}
