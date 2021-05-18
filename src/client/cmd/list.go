package cmd

import (
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

var (
	listDirPath string
	listDirID   uint64
)

// printContents Logs the list directory contents request response in a readable format
func printContents(contents *models.DirectoryContents) {
	if len(contents.Files)+len(contents.SubDirectories) < 1 {
		log.Info("Directory is empty")
		return
	}

	if len(contents.Files) > 0 {
		log.Info("--- Files ---")
	}
	for _, file := range contents.Files {
		log.WithFields(log.Fields{
			"ID":          file.ID,
			"DirectoryID": file.DirectoryID,
			"FileSize":    file.Size,
			"Available":   file.Available,
		}).Infof(color.GreenString(file.Name))
	}
	if len(contents.SubDirectories) > 0 {
		log.Info("--- Directories ---")
	}
	for _, dir := range contents.SubDirectories {
		log.WithFields(log.Fields{"ID": dir.ID, "FileCount": dir.FileCount}).Infof(color.BlueString(dir.Name))
	}
}

// ListDirByPath Sends a request to the server to list the contents of a given directory
// that is specified by the path to the directory
func ListDirByPath() *models.DirectoryContents {
	params := files.NewListDirPathInfoParams().WithDirPath(listDirPath)

	resp, err := apiclient.Client.Files.ListDirPathInfo(params, apiclient.ClientAuth)
	if err != nil {
		log.WithError(err).Fatal("Failed to retrieve directory contents")
	}
	return resp.GetPayload()
}

// ListDirectoryByID Sends a request to the server to list the contents of a given directory
// that is specified by their path from the home dir or their ID
func ListDirectoryByID() *models.DirectoryContents {
	params := files.NewListDirectoryParams().WithID(listDirID)

	resp, err := apiclient.Client.Files.ListDirectory(params, apiclient.ClientAuth)
	if err != nil {
		log.WithError(err).Fatal("Failed to retrieve directory contents")
	}

	log.Infof("Directory ID: %d\n", listDirID)
	return resp.GetPayload()
}

// ListHomeDirectory Sends a request to the server to list the contents of the user's home directory
func ListHomeDirectory() *models.DirectoryContents {
	resp, err := apiclient.Client.Files.ListHomeDirectory(files.NewListHomeDirectoryParams(), apiclient.ClientAuth)
	if err != nil {
		log.WithError(err).Fatal("Failed to retrieve directory contents")
	}

	log.Infof("Home Directory ID: %d\n", resp.Payload.CurrentDir.ID)
	return resp.GetPayload()
}

// listCmd Handles the user inputs from the command line and outputs the result of the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "dir"},
	Short:   "List files on the server",
	Long:    `Returns a list of the files in a specified location on the server`,
	PreRun:  authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		var contents *models.DirectoryContents
		if len(args) > 0 && args[0] != "" {
			listDirPath = args[0]
			contents = ListDirByPath()
		} else if listDirPath != "" {
			contents = ListDirByPath()
		} else if listDirID != 0 {
			contents = ListDirectoryByID()
		} else {
			contents = ListHomeDirectory()
		}
		printContents(contents)
	},
}

func init() {
	listCmd.Flags().StringVarP(&listDirPath, "dir-path", "p", "", "Specify the to a directory to list its contents")
	listCmd.Flags().Uint64VarP(&listDirID, "dir-id", "d", 0, "Specify an ID of a directory to list its contents")
	rootCmd.AddCommand(listCmd)
}
