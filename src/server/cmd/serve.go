package cmd

import (
	"github.com/go-openapi/loads"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Synche API server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		StartSynche(cmd.Flags())
	},
}

func StartSynche(flags *flag.FlagSet) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatal(err)
	}

	port, err := flags.GetInt("port")
	if err != nil {
		log.Fatal(err)
	}
	if port == 0 {
		restapi.OverrideFlags()
	}

	var server *restapi.Server // make sure init is called

	api := operations.NewSyncheAPI(swaggerSpec)
	// get server with flag values filled out
	server = restapi.NewServer(api)
	defer server.Shutdown()

	server.ConfigureAPI()
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}

}

func init() {
	rootCmd.AddCommand(serveCmd)
	if ServerFlags.HasFlags() {
		serveCmd.Flags().AddFlagSet(ServerFlags)
	}

	err := viper.BindPFlag("config.server.port", serveCmd.Flags().Lookup("port"))
	if err != nil {
		log.Fatal(err)
	}
}
