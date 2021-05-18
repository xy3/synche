package repo

import (
	"errors"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
)

// GetHomeDir finds a user's home directory
func GetHomeDir(userID uint, db *gorm.DB) (homeDir *schema.Directory, err error) {
	tx := db.Where("parent_id IS NULL AND user_id = ?", userID).First(&homeDir)
	return homeDir, tx.Error
}

// GetOrCreateHomeDir retrieves the user's home directory, and only creates it on the disk if it is not present
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

// BuildFullPath creates a full path by prepending it with the user's home directory
func BuildFullPath(path string, user *schema.User, db *gorm.DB) (string, error) {
	homeDir, err := GetOrCreateHomeDir(user, db)
	if err != nil {
		return path, err
	}
	return filepath.Join(homeDir.Path, path), nil
}

// GetDirectoryByID retrieves a directory by its ID
func GetDirectoryByID(dirID uint, db *gorm.DB) (dir *schema.Directory, err error) {
	tx := db.First(&dir, dirID)
	return dir, tx.Error
}

// GetDirByPath retrieves a directory by its path
func GetDirByPath(path string, db *gorm.DB) (dir *schema.Directory, err error) {
	path = strings.TrimRight(strings.TrimSpace(path), "/")
	tx := db.Where(schema.Directory{PathHash: hash.PathHash(path)}).First(&dir)
	return dir, tx.Error
}

// GetDirectoryForFileID retrieves the directory a certain file is in
func GetDirectoryForFileID(fileID uint, db *gorm.DB) (*schema.Directory, error) {
	var file schema.File
	res := db.Joins("Directory").First(&file, fileID)
	if res.Error != nil {
		return nil, res.Error
	}

	return file.Directory, nil
}

// GetOrCreateDirectory finds a directory by its path, or creates the it if it does not exist
func GetOrCreateDirectory(path string, db *gorm.DB) (dir *schema.Directory, err error) {
	if dir, err = GetDirByPath(path, db); errors.Is(err, gorm.ErrRecordNotFound) {
		return CreateDirectoryFromPath(path, db)
	}
	return dir, err
}

// CreateDirectoryFromPath uses a directory full-path to create a new directory both on the
// disk and in the database
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

// CreateDirectory creates a new directory inside the parentID directory, both on the disk and in the database
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

// UpdateDirFileCount updates a directorie's file count by querying the database for files
// with the same directory ID
func UpdateDirFileCount(dirID uint, db *gorm.DB) error {
	directory := &schema.Directory{}
	tx := db.Where("id = ?", dirID).First(directory)
	if tx.Error != nil {
		return tx.Error
	}
	_, err := directory.UpdateFileCount(db)
	return err
}

// GenerateUserDirName generates a unique name for a user's home directory using their email
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

// MakeUserHomeDir creates the user's home directory on the disk
func MakeUserHomeDir(user *schema.User) (homeDir string, err error) {
	homeDir = filepath.Join(c.Config.Server.StorageDir, GenerateUserDirName(user))
	err = files.AppFS.MkdirAll(homeDir, 0755)
	return
}

// SetupUserHomeDir creates a user's home directory, both on the disk and in the database
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
