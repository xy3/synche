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
)

func printContents(contents []*models.File) {
	for _, file := range contents {
		fmt.Printf("\nName: \t\t%v\nFile ID: \t%v\nDirectory ID: \t%v\nSize (bytes): \t%v\n", *file.Name, *file.ID, *file.StorageDirectoryID, *file.Size)
	}
}

func listDPathJob(dirPath string) error {
	if err := apiclient.AuthenticateClient(filepath.Join(config.SyncheDir, "token.json")); err != nil {
		return err
	}
	requestAccepted, err := apiclient.Client.Files.ListDPath(files.NewListDPathParams().WithDirPath(dirPath), apiclient.ClientAuth)
	if err != nil {
		return err
	}
	fmt.Printf("Directory path: %v\n", dirPath)
	printContents(requestAccepted.GetPayload().Contents)
	return nil
}

func listDIDJob(dirId uint64) error {
	if err := apiclient.AuthenticateClient(filepath.Join(config.SyncheDir, "token.json")); err != nil {
		return err
	}
	requestAccepted, err := apiclient.Client.Files.ListDID(files.NewListDIDParams().WithDirectoryID(dirId), apiclient.ClientAuth)
	if err != nil {
		return err
	}
	printContents(requestAccepted.GetPayload().Contents)
	return nil
}

var dirPath string
var dirID int64
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files on the server",
	Long:  `Returns a list of the files in a specified location on the server`,
	Run: func(cmd *cobra.Command, args []string) {
		if dirPath != "" {
			err := listDPathJob(dirPath)
			if err != nil {
				log.WithError(err).Fatal("Failed to retrieve directory contents")
			}
		} else if dirID != 0 {
			err := listDIDJob(uint64(dirID))
			if err != nil {
				log.WithError(err).Fatal("Failed to retrieve directory contents")
			}
		} else {
			fmt.Println("Invalid directory specifications. Enter: `synche list -h` for manual")
		}
	},
}

func init() {
	listCmd.Flags().StringVarP(&dirPath, "dir-path", "p", "", "Specify the path to a directory to list its contents. 'home' will list the contents of the storage directory")
	listCmd.Flags().Int64VarP(&dirID, "dir-id", "i", 0, "Specify an ID of a directory to list its contents")
	rootCmd.AddCommand(listCmd)
}
