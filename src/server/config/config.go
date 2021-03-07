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
	Port       string
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

func SetDefaults(home string) interface{} {
	syncheDir := filepath.Join(home, ".synche")
	storageDir := filepath.Join(syncheDir, "data")

	defaultCfg := Configuration{
		Synche: SyncheConfig{
			Dir:     syncheDir,
			Verbose: false,
			Debug:   false,
		},
		Server: ServerConfig{
			Port:       "8080",
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
	err := config.Read(cfgFile, "synche-server", SetDefaults)
	if err != nil {
		return err
	}
	err = viper.UnmarshalKey("config", &Config)
	if err != nil {
		return err
	}

	// Read updates to the config file while server is running
	viper.WatchConfig()
	return nil
}
