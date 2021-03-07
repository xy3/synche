package cmd

import (
	"github.com/go-openapi/loads"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Synche API server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		StartSynche()
	},
}

func StartSynche() {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatal(err)
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
}
