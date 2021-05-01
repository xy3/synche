package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"strconv"
)

func deleteJob(fileId uint64) error {
	_, err := apiclient.Client.Files.DeleteFile(files.NewDeleteFileParams().WithFileID(fileId), apiclient.ClientAuth)
	if err != nil {
		return err
	}

	log.Infof("Deleted file with file ID: %v", fileId)
	return nil
}

var deleteCmd = &cobra.Command{
	Use:    "delete",
	Short:  "Delete a file on the server",
	Long:   `Sends a request to the server to delete file by specified file id.`,
	PreRun: authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		fileId, err := strconv.Atoi(args[0])
		if err != nil {
			log.WithError(err).Fatal("Invalid file id")
		}
		if err = deleteJob(uint64(fileId)); err != nil {
			log.WithError(err).Fatal("Failed to delete file")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
