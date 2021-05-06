package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
)

var (
	deleteFilepath string
	deleteFileID   uint64
)

// deleteJobByPath Sends a request to the server to delete a file specified by an ID
func deleteJobByPath() error {
	_, err := apiclient.Client.Files.DeleteFilepath(files.NewDeleteFilepathParams().WithFilePath(deleteFilepath), apiclient.ClientAuth)
	if err != nil {
		return err
	}
	log.Infof("Deleted file with file ID: %v", deleteFilepath)
	return nil
}

// deleteJobByPath Sends a request to the server to delete a file specified by the path to a file
func deleteJobByID() error {
	_, err := apiclient.Client.Files.DeleteFile(files.NewDeleteFileParams().WithFileID(deleteFileID), apiclient.ClientAuth)
	if err != nil {
		return err
	}
	log.Infof("Deleted file with file ID: %v", deleteFileID)
	return nil
}

// deleteCmd Handles the user inputs from the command line and outputs the result of the delete command
var deleteCmd = &cobra.Command{
	Use:    "delete",
	Short:  "Delete a file on the server",
	Long:   `Sends a request to the server to delete file by specified file id or file path.`,
	PreRun: authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if deleteFilepath != "" {
			if err := deleteJobByPath(); err != nil {
				log.WithError(err).Fatal("Failed to delete file by specified path")
			}
		} else if deleteFileID != 0 {
			if err := deleteJobByID(); err != nil {
				log.WithError(err).Fatal("Failed to delete file by specified ID")
			}
		} else {
			log.Info(cmd.Help())
		}
	},
}

// TODO: fix bug that doesn't allow file paths to be the default args

func init() {
	deleteCmd.Flags().StringVarP(&deleteFilepath, "file-path",
		"p",
		"",
		"Specify the path to a file to delete it")
	deleteCmd.Flags().Uint64VarP(&deleteFileID,
		"file-id",
		"i",
		0,
		"Specify the ID of a file to delete it")
	rootCmd.AddCommand(deleteCmd)
}
