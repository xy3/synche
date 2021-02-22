package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
)

var Config Configuration

type Configuration struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
	UploadDir string
}

type DatabaseConfig struct {
	Driver   string
	Name     string
	Username string
	Password string
	Protocol string
	Address  string
}

func SetDefaults() (err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	syncheDir := path.Join(home, ".synche")

	viper.SetDefault("synche.dir", syncheDir)
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.uploadDir", syncheDir + "/data/received")
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.name", "synche")
	viper.SetDefault("database.username", "admin")
	viper.SetDefault("database.password", "admin")
	viper.SetDefault("database.protocol", "tcp")
	viper.SetDefault("database.address", "127.0.0.1:3306")
	return nil
}

func InitConfig() (err error) {
	err = SetDefaults()

	viper.SetConfigName("synche-server")
	viper.SetConfigType("yaml")

	syncheDir := viper.GetString("synche.dir")
	if _, err := os.Stat(syncheDir); os.IsNotExist(err) {
		err = os.Mkdir(syncheDir, 0755)
	}
	if err != nil {
		return err
	}

	// Set config file
	viper.AddConfigPath(syncheDir)

	// Enable reading environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Config file doesn't exist, create it
		err = viper.WriteConfigAs(path.Join(syncheDir, ".synche-server.yaml"))
		if err != nil {
			return err
		}
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	err = viper.Unmarshal(&Config)

	return err
}
