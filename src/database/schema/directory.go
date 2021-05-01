package schema

import (
	"errors"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gorm.io/gorm"
)

var (
	ErrDeleteHomeDirNotAllowed = errors.New("you cannot delete a home directory")
	ErrDirNotEmpty             = errors.New("directory is not empty")
)

type Directory struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Path      string `gorm:"not null"`
	PathHash  string `gorm:"size:32;uniqueIndex"`
	// Size      int64
	FileCount int64
	UserID    uint
	ParentID  *uint

	Parent   *Directory  `gorm:"foreignKey:id;association_foreignKey:parent_id;association_autoupdate:false;association_autocreate:false"`
	User     User        `gorm:"association_autoupdate:false;association_autocreate:false"`
	Children []Directory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:parent_id;association_autoupdate:false;association_autocreate:false"`
	Files    []File      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:directory_id;association_autoupdate:false;association_autocreate:false"`
}

// executeDelete removes the directory and everything contained within it from the disk
func (dir *Directory) executeDelete() error {
	return files.Afs.RemoveAll(dir.Path)
}

func (dir *Directory) Delete(forceDelete bool, db *gorm.DB) (err error) {
	if err = dir.executeDelete(); err != nil {
		return err
	}

	db = db.Begin()
	if dir.Name == "home" {
		return ErrDeleteHomeDirNotAllowed
	}

	var filesInDir int64
	db.Model(File{}).Where(File{DirectoryID: dir.ID}).Count(&filesInDir)
	if filesInDir != 0 && !forceDelete {
		return ErrDirNotEmpty
	}

	if err = db.Where(File{DirectoryID: dir.ID}).Delete(&File{}).Error; err != nil {
		db.Rollback()
		return err
	}

	if err = db.Unscoped().Delete(&Directory{}, dir).Error; err != nil {
		db.Rollback()
		return err
	}

	db.Commit()
	return nil
}
//
// func (dir *Directory) UpdateSize(size int64, db *gorm.DB) (err error) {
// 	dir.Size += size
//
// 	if dir.Parent == nil {
// 		db.Preload("Parent").Find(dir)
// 	}
//
// 	if err = db.Save(dir).Error; err != nil {
// 		return err
// 	}
//
// 	if dir.Parent == nil {
// 		return nil
// 	}
//
// 	// Recursively update parent directories
// 	return dir.Parent.UpdateSize(size, db)
// }

func (dir *Directory) UpdateFileCount(db *gorm.DB) (num int64, err error) {
	tx := db.Model(&File{}).Where("directory_id = ?", dir.ID).Count(&num)
	if tx.Error != nil {
		return num, tx.Error
	}
	dir.FileCount = num
	return num, db.Save(dir).Error
}
