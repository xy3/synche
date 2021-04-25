package cmd

import (
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
	"path/filepath"
)

func printContents(contents *models.DirectoryContents) {
	if len(contents.Files)+len(contents.Subdirectories) < 1 {
		log.Info("Directory is empty")
		return
	}

	if len(contents.Files) > 0 {
		log.Info("--- Files ---")
	}
	for _, file := range contents.Files {
		log.WithFields(log.Fields{
			"ID":          *file.ID,
			"DirectoryID": *file.StorageDirectoryID,
			"FileSize":    *file.Size,
		}).Infof(color.GreenString(*file.Name))
	}
	if len(contents.Subdirectories) > 0 {
		log.Info("--- Directories ---")
	}
	for _, dir := range contents.Subdirectories {
		log.WithFields(log.Fields{"ID": dir.ID, "FileCount": dir.FileCount}).Infof(color.BlueString(dir.Name))
	}
}

func listDirectoryJob(params *files.ListDirectoryParams) *files.ListDirectoryOK {
	resp, err := apiclient.Client.Files.ListDirectory(params, apiclient.ClientAuth)
	if err != nil {
		log.WithError(err).Fatal("Failed to retrieve directory contents")
	}
	return resp
}

func ListDirectoryByName(name string) *models.DirectoryContents {
	resp := listDirectoryJob(files.NewListDirectoryParams().WithDirectoryName(&name))
	log.Infof("Directory name: %v\n", name)
	return resp.GetPayload()
}

func ListDirectoryByID(dirId uint64) *models.DirectoryContents {
	resp := listDirectoryJob(files.NewListDirectoryParams().WithDirectoryID(&dirId))
	log.Infof("Directory ID: %d\n", dirId)
	return resp.GetPayload()
}

var listDirName string
var listDirID int64
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files on the server",
	Long:  `Returns a list of the files in a specified location on the server`,
	Run: func(cmd *cobra.Command, args []string) {
		err := apiclient.Authenticator(filepath.Join(config.SyncheDir, "token.json"))
		if err != nil {
			log.Fatal("Failed to authenticate the client")
		}

		var contents *models.DirectoryContents
		if listDirName != "" {
			contents = ListDirectoryByName(listDirName)
		} else if listDirID != 0 {
			contents = ListDirectoryByID(uint64(listDirID))
		} else {
			contents = ListDirectoryByName("home")
		}
		printContents(contents)
	},
}

func init() {
	listCmd.Flags().StringVarP(&listDirName, "dir-name", "n", "", "Specify the name to a directory to list its contents. Defaults to your base user directory")
	listCmd.Flags().Int64VarP(&listDirID, "dir-id", "i", 0, "Specify an ID of a directory to list its contents")
	rootCmd.AddCommand(listCmd)
}
