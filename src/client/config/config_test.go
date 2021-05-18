package config

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
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