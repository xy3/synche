package setup_test

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/setup"
	"testing"
)

func TestDirs(t *testing.T) {
	files.SetFileSystem(afero.NewMemMapFs())
	dirs := []string{"dir1", "dir2", "/nested/dir", "relative/dir"}
	err := setup.Dirs(files.AppFS, dirs)
	assert.NoError(t, err)
	for _, dir := range dirs {
		isDir, _ := files.Afs.IsDir(dir)
		assert.True(t, isDir)
	}
}
