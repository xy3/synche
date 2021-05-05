package repo

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
)

func GetHomeDir(userID uint) (*schema.Directory, error) {
	homeDir := &schema.Directory{}
	tx := database.DB.Where("parent_id IS NULL AND user_id = ?", userID).First(homeDir)
	return homeDir, tx.Error
}

func GetDirectoryByID(dirID uint) (*schema.Directory, error) {
	var directory schema.Directory
	if err := database.DB.First(&directory, dirID).Error; err != nil {
		return nil, err
	}
	return &directory, nil
}

func GetDirectoryForFileID(fileId uint) (*schema.Directory, error) {
	var file schema.File
	res := database.DB.Joins("Directory").Find(&file, fileId)
	if res.Error != nil {
		return nil, res.Error
	}

	return file.Directory, nil
}

func GetDirByPath(path string) (*schema.Directory, error) {
	log.Infof("received GetDirByPath request for: %s", path)
	pathHash := hash.PathHash(path)
	dir := &schema.Directory{}
	tx := database.DB.Where(&schema.Directory{PathHash: pathHash}).First(dir)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return dir, nil
}

func CreateDirectory(name string, parentID uint, db *gorm.DB) (directory *schema.Directory, err error) {
	name = strings.Trim(strings.TrimSpace(name), "/")
	if name == "" {
		return nil, errors.New("directory name is invalid")
	}

	parent := schema.Directory{}
	if err = db.Where("id = ?", parentID).First(&parent).Error; err != nil {
		return nil, err
	}

	newPath := filepath.Join(parent.Path, name)
	directory = &schema.Directory{
		Name:     name,
		Path:     newPath,
		PathHash: hash.PathHash(newPath),
		UserID:   parent.UserID,
		ParentID: &parent.ID,
	}

	if err = files.Afs.MkdirAll(directory.Path, 0755); err != nil {
		return nil, err
	}

	db = db.Begin()
	if err = db.Create(directory).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	db.Commit()

	return directory, nil
}

func CreateDirectoryFromPath(path string, db *gorm.DB) (dir *schema.Directory, err error) {
	newPath := strings.TrimRight(strings.TrimSpace(path), "/")
	parts := strings.Split(newPath, "/")

	parentPath := filepath.Dir(newPath)
	parentPathHash := hash.PathHash(parentPath)
	parentDir := &schema.Directory{}
	tx := db.Where(schema.Directory{PathHash: parentPathHash}).First(parentDir)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		if len(parts) < len(strings.Split(c.Config.Server.StorageDir, "/")) {
			return nil, tx.Error
		}
		return CreateDirectoryFromPath(parentPath, db)
	}

	return CreateDirectory(parts[len(parts)-1], parentDir.ID, db)
}

func UpdateDirFileCount(dirID uint) error {
	directory := &schema.Directory{}
	tx := database.DB.Where("id = ?", dirID).First(directory)
	if tx.Error != nil {
		return tx.Error
	}
	_, err := directory.UpdateFileCount(database.DB)
	return err
}

func GenerateUserDirName(user *schema.User) string {
	var userSlug = user.Email
	userSlug = strings.ReplaceAll(userSlug, "@", "")
	userSlug = strings.ReplaceAll(userSlug, ".", "")
	userSlug = userSlug[:5]
	userHash := hash.MD5Hash([]byte(user.Email + user.Password))
	return userSlug + "_" + userHash
}

func CreateUserHomeDir(user *schema.User) (homeDir string, err error) {
	homeDir = filepath.Join(c.Config.Server.StorageDir, GenerateUserDirName(user))
	err = files.AppFS.MkdirAll(homeDir, 0755)
	return
}

func SetupUserHomeDir(user *schema.User) (*schema.Directory, error) {
	// Create the user's home directory
	homeDirPath, err := CreateUserHomeDir(user)
	if err != nil {
		return nil, err
	}

	homeDir := &schema.Directory{
		Name:     "home",
		Path:     homeDirPath,
		PathHash: hash.PathHash(homeDirPath),
		UserID:   user.ID,
		ParentID: nil,
	}

	if err = database.DB.Create(homeDir).Error; err != nil {
		return nil, err
	}

	return homeDir, nil
}

func GetTotalFileChunks(fileID uint64) (uint64, error) {
	var file schema.File
	tx := database.DB.Where("file.id = ?", fileID).First(&file)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return uint64(file.TotalChunks), nil
}
