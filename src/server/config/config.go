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

type Server struct {
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

func SetDefaults() (err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	syncheDir := path.Join(home, ".synche")

	viper.SetDefault("synche.dir", syncheDir)
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.uploadDirectory", syncheDir + "/data/received")
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

	viper.SetConfigName(".synche-server")
	viper.SetConfigType("yaml")

	syncheDir := viper.GetString("synche.dir")
	if _, err := os.Stat(syncheDir); os.IsNotExist(err) {
		err = os.Mkdir(syncheDir, os.ModePerm)}
	if err != nil {
		return err
	}

	// Set config file
	viper.AddConfigPath(syncheDir)
	cfgFile := path.Join(syncheDir, ".synche-server.yaml")
	viper.SetConfigFile(cfgFile)

	// Enable reading environment variables
	viper.AutomaticEnv()

	// If unmarshaling, create configuration here
	// var configuration Configuration

	if err := viper.ReadInConfig(); err != nil {
		// Config file doesn't exist, create it
		err = viper.WriteConfigAs(cfgFile)
		if err != nil {
			return err
		}
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	// Unmarshalling can be done here
	// err = viper.Unmarshal(&configuration)

	return err
}
