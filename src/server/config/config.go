package config

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
	"path/filepath"
)

var Config = &Configuration{
	Synche: SyncheConfig{
		Dir:     config.SyncheDir,
		Verbose: false,
		Debug:   false,
	},
	Server: ServerConfig{
		Port:       9449,
		Host:       "127.0.0.1",
		BasePath:   "/v1/api",
		StorageDir: filepath.Join(config.SyncheDir, "data"),
		UploadDir:  filepath.Join(config.SyncheDir, "data", "received"),
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
	Host       string
	BasePath   string
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

func init() {
	viper.SetDefault("config", Config)
}

func InitConfig(cfgFile string) error {
	_, err := config.ReadOrCreate("synche-server", cfgFile, Config, Config)
	if err != nil {
		return err
	}

	err = viper.UnmarshalKey("config", &Config)
	if err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Info("Config Updated")
		_ = viper.UnmarshalKey("config", &Config)
	})
	return nil
}
