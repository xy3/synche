package jobs

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"path/filepath"
	"strings"
)

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

func CreateDirectoryInsideUserHomeDir(name string, user *schema.User) (string, error) {
	homeDir := filepath.Join(c.Config.Server.StorageDir, GenerateUserDirName(user))
	newDir := filepath.Join(homeDir, name)
	return newDir, files.AppFS.MkdirAll(newDir, 0755)
}

func DeleteDirectory(path string) error {
	return files.AppFS.RemoveAll(path)
}