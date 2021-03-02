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
	Redis RedisConfig
}

type ServerConfig struct {
	Port string
	UploadDir string
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
	Network string
	Address string
	Port string
	Password string
	DB int
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
	viper.SetDefault("server.storageDir", home)

	viper.SetDefault("data.driver", "mysql")
	viper.SetDefault("data.name", "synche")
	viper.SetDefault("data.username", "admin")
	viper.SetDefault("data.password", "admin")
	viper.SetDefault("data.protocol", "tcp")
	viper.SetDefault("data.address", "127.0.0.1:3306")

	viper.SetDefault("redis.network", "tcp")
	viper.SetDefault("redis.address", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
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
		err = viper.SafeWriteConfigAs(path.Join(syncheDir, "synche-server.yaml"))
		if err != nil {
			return err
		}
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	err = viper.Unmarshal(&Config)

	return err
}
