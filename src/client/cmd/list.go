package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
	"path/filepath"
	"strconv"
)

func getFileNames(dirContents []*models.File) []string {
	var fileNames []string
	for _, file := range dirContents {
		fileNames = append(fileNames, *file.Name)
	}
	return fileNames
}

func listJob(dirId uint64) error {
	if err := apiclient.AuthenticateClient(filepath.Join(config.SyncheDir, "token.json")); err != nil {
		return err
	}
	requestAccepted, err := apiclient.Client.Files.List(files.NewListParams().WithDirectoryID(dirId), apiclient.ClientAuth)
	if err != nil {
		return err
	}
	fmt.Printf("Directory contents: %v\n", getFileNames(requestAccepted.GetPayload().Contents))
	return nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files on the server",
	Long:  `Returns a list of the files in a specified location on the server`,
	Run: func(cmd *cobra.Command, args []string) {
		dirId, err := strconv.Atoi(args[0])
		if err != nil {
			log.WithError(err).Fatal("Invalid directory id")
		}
		if err = listJob(uint64(dirId)); err != nil {
			log.WithError(err).Fatal("Failed to retrieve directory contents")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
