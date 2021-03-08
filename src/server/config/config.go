package config

import (
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
	"path/filepath"
)

var Config Configuration

type Configuration struct {
	Synche   SyncheConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

type SyncheConfig struct {
	Dir     string
	Verbose bool
	Debug   bool
}

type ServerConfig struct {
	Port       int
	UploadDir  string
	StorageDir string
}

type DatabaseConfig struct {
	Driver   string
	Name     string
	Username string
	Password string
	Protocol string
	Address  string
}

type RedisConfig struct {
	Protocol string
	Address  string
	Password string
	DB       int
}

func RequiredDirs() []string {
	return []string{
		Config.Synche.Dir,
		Config.Server.UploadDir,
		Config.Server.StorageDir,
	}
}

func ServerDefaults(syncheDir string) interface{} {
	storageDir := filepath.Join(syncheDir, "data")

	defaultCfg := Configuration{
		Synche: SyncheConfig{
			Dir:     syncheDir,
			Verbose: false,
			Debug:   false,
		},
		Server: ServerConfig{
			Port:       8080,
			StorageDir: storageDir,
			UploadDir:  filepath.Join(storageDir, "received"),
		},
		Database: DatabaseConfig{
			Driver:   "mysql",
			Name:     "synche",
			Username: "root",
			Password: "",
			Protocol: "tcp",
			Address:  "127.0.0.1:3306",
		},
		Redis: RedisConfig{
			Protocol: "tcp",
			Address:  "127.0.0.1:6379",
			Password: "",
			DB:       0,
		},
	}

	viper.SetDefault("config", defaultCfg)
	return defaultCfg
}

func InitConfig(cfgFile string) error {
	cfg, err := config.New(cfgFile, "synche-server")
	if err != nil {
		return err
	}

	err = cfg.ReadOrCreate(ServerDefaults(cfg.Dir))
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
