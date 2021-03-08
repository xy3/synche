package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/setup"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "synche",
	Short: "Quickly upload and manage your files on a Synche server",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.OnInitialize(func() {
		err := c.InitConfig(cfgFile)
		if err != nil {
			log.WithError(err).Fatal("Could not initialize the config")
		}

		err = setup.Dirs(files.AppFS, c.RequiredDirs())
		if err != nil {
			log.WithError(err).Fatal("Could not set up the required directories")
		}

		err = c.ConfigureClient()
		if err != nil {
			log.WithError(err).Fatal("Failed to configure the Synche client")
		}

		if c.Config.Synche.Debug {
			log.Infof("Verbose: true")
			log.SetLevel(log.DebugLevel)
		}
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.synche/synche-client.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "display verbose output (default is false)")
	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		log.Fatal(err)
	}
}
