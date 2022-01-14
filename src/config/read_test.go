package config_test

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/xy3/synche/src/config"
	"github.com/xy3/synche/src/files"
	c "github.com/xy3/synche/src/server"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	files.SetFileSystem(afero.NewMemMapFs())
	viper.SetFs(files.AppFS)
	os.Exit(m.Run())
}

func TestRead(t *testing.T) {
	t.Run("Existing config", func(t *testing.T) {
		err := config.Write("test-config.yml", c.Config)
		assert.NoError(t, err)

		err = config.Read("test-config", "test-config.yml")
		assert.NoError(t, err)
		assert.Equal(t, "test-config.yml", viper.ConfigFileUsed())
		assert.True(t, viper.IsSet("config"))

		err = viper.UnmarshalKey("config", &c.Config)
		assert.NoError(t, err)
	})

	t.Run("Non-existing config", func(t *testing.T) {
		err := config.Read("test-config2", "")
		assert.Error(t, err)
	})
}

func TestReadOrCreate(t *testing.T) {
	config.TestMode = true

	t.Run("Existing config", func(t *testing.T) {
		created, err := config.ReadOrCreate("test-config", "test-config.yml", c.Config, c.Config)
		assert.NoError(t, err)
		assert.False(t, created)

		assert.Equal(t, "test-config.yml", viper.ConfigFileUsed())
		assert.True(t, viper.IsSet("config"))

		err = viper.UnmarshalKey("config", &c.Config)
		assert.NoError(t, err)
	})

	t.Run("Non-existing config", func(t *testing.T) {
		created, err := config.ReadOrCreate("test-config2", "", c.Config, c.Config)
		assert.NoError(t, err)
		assert.True(t, created)

		assert.Equal(t, filepath.Join(config.SyncheDir, "test-config2.yaml"), viper.ConfigFileUsed())
		assert.True(t, viper.IsSet("config"))

		err = viper.UnmarshalKey("config", &c.Config)
		assert.NoError(t, err)
	})
}
