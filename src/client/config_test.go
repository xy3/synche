package client

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/xy3/synche/src/config"
	"github.com/xy3/synche/src/files"
	"testing"
)

func TestMain(m *testing.M) {
	config.TestMode = true
	files.SetFileSystem(afero.NewMemMapFs())
	viper.SetFs(files.AppFS)
}

func TestRequiredDirs(t *testing.T) {
	want := []string{
		Config.Synche.Dir,
		Config.Synche.DataDir,
	}
	got := RequiredDirs()
	assert.Equal(t, want, got)
}
func TestInitConfig(t *testing.T) {
	err := InitConfig("")
	assert.NoError(t, err)
}
