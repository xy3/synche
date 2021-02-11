package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	viper.AddConfigPath(home)

	// Enable reading environment variables
	viper.AutomaticEnv()

	// If unmarshaling, create configuration here
	// var configuration Configuration

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Printf("Unable to decode into struct, %v\n", err)
	}

	// err = viper.Unmarshal(&configuration)

	return err
}
