package client

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/xy3/synche/src/config"
	"path"
)

var Config = Configuration{
	Synche: SyncheConfig{
		Dir:     config.SyncheDir,
		DataDir: path.Join(config.SyncheDir, "data"),
	},
	Chunks: ChunksConfig{
		SizeKB:  1024,
		Workers: 10,
	},
	Server: ServerConfig{
		Host:     DefaultHost,
		BasePath: DefaultBasePath,
		Https:    len(DefaultSchemes) > 1,
	},
}

type Configuration struct {
	Synche SyncheConfig
	Chunks ChunksConfig
	Server ServerConfig
}

type ServerConfig struct {
	Host     string
	BasePath string
	Https    bool
}

type SyncheConfig struct {
	Dir     string
	DataDir string
	Verbose bool
	Debug   bool
}

type ChunksConfig struct {
	SizeKB  int64
	Workers int
}

func RequiredDirs() []string {
	return []string{
		Config.Synche.Dir,
		Config.Synche.DataDir,
	}
}

func init() {
	viper.SetDefault("config", &Config)
}

func InitConfig(cfgFile string) error {
	_, err := config.ReadOrCreate("synche-client", cfgFile, Config, Config)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("config", &Config)
	if err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		_ = viper.UnmarshalKey("config", &Config)
	})
	return nil
}
