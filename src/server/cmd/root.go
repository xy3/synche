package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/setup"
)

var (
	ServerFlags *flag.FlagSet
	cfgFile     string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your application",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.OnInitialize(func() {
		if viper.GetBool("config.synche.debug") {
			log.Infof("Debug: true")
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

	err = c.InitConfig(cfgFile)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize config")
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", c.Config.Synche.Dir, "config file")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "display debug output")
	err = viper.BindPFlag("config.synche.debug", rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		log.Fatal(err)
	}
}
