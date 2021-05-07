package config

import (
	"fmt"
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
		StorageDir: filepath.Join(config.SyncheDir, "data", "storage"),
		ChunkDir:   filepath.Join(config.SyncheDir, "data", "chunks"),
		SecretKey:  "CHANGE_THIS_SECRET_KEY",
	},
	Database: DatabaseConfig{
		Driver:   "mysql",
		Name:     "synche",
		Username: "root",
		Password: "",
		Protocol: "tcp",
		Address:  "127.0.0.1:3306",
	},
	Ftp: FtpConfig{
		Port:           2121,
		KeyFile:        "",
		Hostname:       "127.0.0.1",
		CertFile:       "",
		PublicIp:       "",
		PassivePorts:   "52013-52114",
		WelcomeMessage: "Welcome to the Synche FTP server",
	},
}

type Configuration struct {
	Synche   SyncheConfig
	Server   ServerConfig
	Database DatabaseConfig
	Ftp      FtpConfig
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
	ChunkDir   string
	StorageDir string
	SecretKey  string
}

type DatabaseConfig struct {
	Driver   string
	Name     string
	Username string
	Password string
	Protocol string
	Address  string
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Protocol,
		c.Address,
		c.Name,
	)
}

type FtpConfig struct {
	Port           int
	KeyFile        string
	Hostname       string
	CertFile       string
	PublicIp       string
	PassivePorts   string
	WelcomeMessage string
}

func RequiredDirs() []string {
	return []string{
		Config.Synche.Dir,
		Config.Server.ChunkDir,
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
