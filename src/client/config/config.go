package config

import (
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
	"path"
)

var (
	Config Configuration
)

type Configuration struct {
	Synche  SyncheConfig
	Chunks  ChunksConfig
	Server  ServerConfig
}

type ServerConfig struct {
	Host     string
	BasePath string
	Schemes  []string
}

type SyncheConfig struct {
	Dir     string
	DataDir string
	Verbose bool
	Debug   bool
}

type ChunksConfig struct {
	Size int64
}

func RequiredDirs() []string {
	return []string{
		Config.Synche.Dir,
		Config.Synche.DataDir,
	}
}

func SetClientDefaults(home string) interface{} {
	syncheDir := path.Join(home, ".synche")
	dataDir := path.Join(syncheDir, "data")

	defaultCfg := Configuration{
		Synche: SyncheConfig{
			Dir:     syncheDir,
			DataDir: dataDir,
			Verbose: false,
			Debug:   false,
		},
		Chunks: ChunksConfig{
			Size: 1,
		},
		Server: ServerConfig{
			Host:     apiclient.DefaultHost,
			BasePath: apiclient.DefaultBasePath,
			Schemes:  apiclient.DefaultSchemes,
		},
	}

	viper.SetDefault("config", defaultCfg)
	return defaultCfg
}

func InitConfig(cfgFile string) error {
	err := config.Read(cfgFile, "synche-client", SetClientDefaults)
	if err != nil {
		return err
	}
	err = viper.UnmarshalKey("config", &Config)
	if err != nil {
		return err
	}

	return nil
}