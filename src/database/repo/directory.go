package repo

import (
	"errors"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
)

func GetHomeDir(userID uint, db *gorm.DB) (homeDir *schema.Directory, err error) {
	tx := db.Where("parent_id IS NULL AND user_id = ?", userID).First(&homeDir)
	return homeDir, tx.Error
}

func GetOrCreateHomeDir(user *schema.User, db *gorm.DB) (homeDir *schema.Directory, err error) {
	homeDir, err = GetHomeDir(user.ID, db)
	if err != nil {
		return
	}

	if exists, _ := files.Afs.IsDir(homeDir.Path); !exists {
		if _, err = MakeUserHomeDir(user); err != nil {
			return
		}
	}

	return homeDir, nil
}

func BuildFullPath(path string, user *schema.User, db *gorm.DB) (string, error) {
	homeDir, err := GetOrCreateHomeDir(user, db)
	if err != nil {
		return path, err
	}
	return filepath.Join(homeDir.Path, path), nil
}

func GetDirectoryByID(dirID uint, db *gorm.DB) (dir *schema.Directory, err error) {
	tx := db.First(&dir, dirID)
	return dir, tx.Error
}

func GetDirectoryForFileID(fileId uint, db *gorm.DB) (*schema.Directory, error) {
	var file schema.File
	res := db.Joins("Directory").Find(&file, fileId)
	if res.Error != nil {
		return nil, res.Error
	}

	return file.Directory, nil
}

func GetDirByPath(path string, db *gorm.DB) (dir *schema.Directory, err error) {
	path = strings.TrimRight(strings.TrimSpace(path), "/")
	tx := db.Where(schema.Directory{PathHash: hash.PathHash(path)}).First(&dir)
	return dir, tx.Error
}

func GetOrCreateDirectory(path string, db *gorm.DB) (dir *schema.Directory, err error) {
	if dir, err = GetDirByPath(path, db); errors.Is(err, gorm.ErrRecordNotFound) {
		return CreateDirectoryFromPath(path, db)
	}
	return dir, err
}

func CreateDirectoryFromPath(path string, db *gorm.DB) (dir *schema.Directory, err error) {
	path = strings.TrimRight(strings.TrimSpace(path), "/")
	parts := strings.Split(path, "/")

	if parts[0] != "" {
		return nil, errors.New("no suitable parent directory could be found")
	}

	parentPath := filepath.Dir(path)

	parentDir, err := GetOrCreateDirectory(parentPath, db)
	if err != nil {
		return nil, err
	}

	return CreateDirectory(parts[len(parts)-1], parentDir.ID, db)
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

	if err = db.Create(directory).Error; err != nil {
		return nil, err
	}

	return directory, nil
}

func UpdateDirFileCount(dirID uint, db *gorm.DB) error {
	directory := &schema.Directory{}
	tx := db.Where("id = ?", dirID).First(directory)
	if tx.Error != nil {
		return tx.Error
	}
	_, err := directory.UpdateFileCount(database.DB)
	return err
}

func GenerateUserDirName(user *schema.User) string {
	var userSlug = user.Email
	userSlug = strings.Split(userSlug, "@")[0]
	userSlug = strings.ReplaceAll(userSlug, ".", "")
	length := len(userSlug)
	if length > 20 {
		length = 20
	}
	userSlug = userSlug[:length]
	userHash := hash.MD5HashString(user.Email + user.Name)[:10]
	return userSlug + "_" + userHash
}

func MakeUserHomeDir(user *schema.User) (homeDir string, err error) {
	homeDir = filepath.Join(c.Config.Server.StorageDir, GenerateUserDirName(user))
	err = files.AppFS.MkdirAll(homeDir, 0755)
	return
}

func SetupUserHomeDir(user *schema.User, db *gorm.DB) (*schema.Directory, error) {
	// Create the user's home directory
	homeDirPath, err := MakeUserHomeDir(user)
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

	if err = db.Create(homeDir).Error; err != nil {
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
