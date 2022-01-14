package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xy3/synche/src/client"
	"github.com/xy3/synche/src/config"
	"github.com/xy3/synche/src/files"
	"path/filepath"
)

var cfgFile string

// rootCmd Represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "synche",
	Short: "Quickly upload and manage your files on a Synche server",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute Adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	InitClientCLI()

	cobra.OnInitialize(func() {
		client.ConfigureClient(client.Config.Server.Host, client.Config.Server.BasePath)

		if client.Config.Synche.Debug {
			log.Debug("Debug mode")
			log.SetLevel(log.DebugLevel)
		}
	})

	cobra.CheckErr(rootCmd.Execute())
}

// authenticateUserPreRun Ensures the user is authenticated
func authenticateUserPreRun(*cobra.Command, []string) {
	err := client.Authenticator(filepath.Join(config.SyncheDir, "token.json"))
	if err != nil {
		log.Fatal("Failed to authenticate the client")
	}
}

func InitClientCLI() {
	err := files.SetupDirs(files.AppFS, client.RequiredDirs())
	if err != nil {
		log.WithError(err).Fatal("Could not set up the required directories")
	}
	// The config needs to be initialized here so that sub-command flags can override the values
	err = client.InitConfig(cfgFile)
	if err != nil {
		log.WithError(err).Fatal("Could not initialize the config")
	}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.synche/synche-client.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&client.Config.Synche.Verbose, "verbose", "v", false, "display verbose output (default is false)")
}
