package cmd

import (
	"context"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
	"path/filepath"

	"github.com/spf13/cobra"
)

// mkdirCmd represents the mkdir command
var mkdirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createDirJob(args[0])
	},
}

func init() {
	rootCmd.AddCommand(mkdirCmd)
}

func createDirJob(name string) {
	if len(name) < 3 {
		log.Error("directory name must be more than 3 characters long")
		return
	}

	err := apiclient.Authenticator(filepath.Join(config.SyncheDir, "token.json"))
	if err != nil {
		log.WithError(err).Fatal("Failed to authenticate the client")
	}

	directory, err := apiclient.Client.Files.CreateDirectory(&files.CreateDirectoryParams{
		DirectoryName: name,
		Context:       context.Background(),
	}, apiclient.ClientAuth)

	if err != nil {
		log.WithError(err).Error("failed to create the directory")
		return
	}

	log.Info("Created the directory successfully.")
	log.Infof("Directory ID: %d", directory.Payload.ID)
}