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
	Synche SyncheConfig
	Chunks ChunksConfig
	Server ServerConfig
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
	Size    int64
	Workers int
}

func RequiredDirs() []string {
	return []string{
		Config.Synche.Dir,
		Config.Synche.DataDir,
	}
}

func ClientDefaults(syncheDir string) interface{} {
	dataDir := path.Join(syncheDir, "data")

	defaultCfg := Configuration{
		Synche: SyncheConfig{
			Dir:     syncheDir,
			DataDir: dataDir,
			Verbose: false,
			Debug:   false,
		},
		Chunks: ChunksConfig{
			Size:    1,
			Workers: 10,
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
	cfg, err := config.New(cfgFile, "synche-client")
	if err != nil {
		return err
	}

	err = cfg.ReadOrCreate(ClientDefaults(cfg.Dir))
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("config", &Config)
	if err != nil {
		return err
	}

	viper.WatchConfig()
	return nil
}
