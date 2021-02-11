package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
)

var (
	cfgFile string
)

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
		err := config.InitConfig(cfgFile)
		if err != nil {
			log.Fatalf("Could not initialize the config: %v", err)
		}
		data.SetupDirs()

		if viper.GetBool("verbose") {
			log.Infof("Verbose: true")
			log.SetLevel(log.DebugLevel)
		}
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.synche.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "display verbose output (default is false)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		panic(err) // TODO
	}
}
