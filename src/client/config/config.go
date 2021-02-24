package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
)

var Config Configuration

type Configuration struct {
	Synche  SyncheConfig
	Chunks  ChunksConfig
	Verbose bool
}

type SyncheConfig struct {
	Dir     string
	DataDir string
}

type ChunksConfig struct {
	Size int64
}

func SetDefaults() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	syncheDir := path.Join(home, ".synche")
	dataDir := path.Join(syncheDir, "data")
	viper.SetDefault("synche.dir", syncheDir)
	viper.SetDefault("synche.dataDir", dataDir)
	viper.SetDefault("chunks.size", 1) // 1MB
	viper.SetDefault("verbose", false)
	return nil
}

// initConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string) error {
	err := SetDefaults()
	if err != nil {
		log.Errorf("Failed to set config defaults: %v", err)
		return err
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("synche-client")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// cfgFile not set, scan usual directories for existing config
		viper.AddConfigPath(viper.GetString("synche.dir"))
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	} else {
		// the config file does not exist, so create a new one
		err = viper.WriteConfigAs(path.Join(viper.GetString("synche.dir"), "synche-client.yaml"))
		if err != nil {
			log.Fatalf("Unable to create new config file, %v", err)
		}
	}

	err = viper.Unmarshal(&Config)

	return nil
}
