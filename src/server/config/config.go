package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
)

// TODO: decide whether to create these structs and unmarshal
/* type Configuration struct {
	Server ServerConfigurations
	Database DatabaseConfigurations
}

type ServerConfigurations struct {
	Port string
}

type DatabaseConfigurations struct {
	DBDriver string
	DBName string
	DBUsername string
	DBPassword string
	DBProtocol string
	DBAddress string
} */

func InitConfig() (err error) {
	viper.SetConfigName(".synche-server")
	viper.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.AddConfigPath(home)
	cfgFile := path.Join(home, ".synche-server.yaml")
	viper.SetConfigFile(cfgFile)


	// Enable reading environment variables
	viper.AutomaticEnv()

	// If unmarshaling, create configuration here
	// var configuration Configuration

	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file:", viper.ConfigFileUsed())
	} else {
		return err
	}

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.uploadDirectory", "../data/received")
	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.name", "synche")
	viper.SetDefault("database.username", "admin")
	viper.SetDefault("database.password", "admin")
	viper.SetDefault("database.protocol", "tcp")
	viper.SetDefault("database.address", "127.0.0.1:3306")

	// err = viper.Unmarshal(&configuration)

	return err
}
