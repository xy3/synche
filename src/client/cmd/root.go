package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/setup"
	"path/filepath"
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
		apiclient.ConfigureClient(c.Config.Server.Host, c.Config.Server.BasePath)

		if c.Config.Synche.Debug {
			log.Debug("Debug mode")
			log.SetLevel(log.DebugLevel)
		}
	})

	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	err := setup.Dirs(files.AppFS, c.RequiredDirs())
	if err != nil {
		log.WithError(err).Fatal("Could not set up the required directories")
	}
	// The config needs to be initialized here so that sub-command flags can override the values
	err = c.InitConfig(cfgFile)
	if err != nil {
		log.WithError(err).Fatal("Could not initialize the config")
	}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.synche/synche-client.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&c.Config.Synche.Verbose, "verbose", "v", false, "display verbose output (default is false)")
}

func authenticateUserPreRun(*cobra.Command, []string) {
	err := apiclient.Authenticator(filepath.Join(config.SyncheDir, "token.json"))
	if err != nil {
		log.Fatal("Failed to authenticate the client")
	}
}
