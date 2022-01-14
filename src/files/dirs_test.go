package files_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/xy3/synche/src/files"
	"testing"
)

func TestDirs(t *testing.T) {
	files.SetFileSystem(afero.NewMemMapFs())
	dirs := []string{"dir1", "dir2", "/nested/dir", "relative/dir"}
	err := files.SetupDirs(files.AppFS, dirs)
	assert.NoError(t, err)
	for _, dir := range dirs {
		isDir, _ := files.Afs.IsDir(dir)
		assert.True(t, isDir)
	}
}
