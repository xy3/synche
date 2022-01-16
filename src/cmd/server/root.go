package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/xy3/synche/src/files"
	"github.com/xy3/synche/src/server"
)

var (
	ServerFlags *flag.FlagSet
	cfgFile     string
)

// rootCmd Represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "The Synche Server command line interface",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	InitServerCLI()

	cobra.OnInitialize(func() {
		if viper.GetBool("config.synche.debug") {
			log.Infof("Debug: true")
			log.SetLevel(log.DebugLevel)
		}
	})

	cobra.CheckErr(rootCmd.Execute())
}

func InitServerCLI() {
	err := files.SetupDirs(files.AppFS, server.RequiredDirs())
	if err != nil {
		log.WithError(err).Fatal("Could not set up the required directories")
	}

	err = server.InitConfig(cfgFile)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize config")
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", server.Config.Synche.Dir, "config file")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "display debug output")
	err = viper.BindPFlag("config.synche.debug", rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		log.Fatal(err)
	}
}
