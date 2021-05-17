package cmd

import (
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

// var listDirName string
var listDirID int64

// TODO Fix bug to allow file paths to be the default arguments

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

// TODO: Replace this with a function to list by path
// func ListDirectoryByName(name string) *models.DirectoryContents {
// 	resp := listDirectoryJob(files.NewListDirectoryParams().WithDirectoryName(&name))
// 	log.Infof("Directory name: %v\n", name)
// 	return resp.GetPayload()
// }

// ListDirectoryByID Sends a request to the server to list the contents of a given directory
// that is specified by their path from the home dir or their ID
func ListDirectoryByID(dirId uint64) *models.DirectoryContents {
	params := files.NewListDirectoryParams().WithID(dirId)

	resp, err := apiclient.Client.Files.ListDirectory(params, apiclient.ClientAuth)
	if err != nil {
		log.WithError(err).Fatal("Failed to retrieve directory contents")
	}

	log.Infof("Directory ID: %d\n", dirId)
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
		// if listDirName != "" {
		// 	contents = ListDirectoryByName(listDirName)
		var contents *models.DirectoryContents
		if listDirID != 0 {
			contents = ListDirectoryByID(uint64(listDirID))
		} else {
			contents = ListHomeDirectory()
		}
		printContents(contents)
	},
}

func init() {
	// listCmd.Flags().StringVarP(&listDirName, "dir-name", "n", "", "Specify the name to a directory to list its contents. Defaults to your base user directory")
	listCmd.Flags().Int64VarP(&listDirID, "dir-id", "d", 0, "Specify an ID of a directory to list its contents")
	rootCmd.AddCommand(listCmd)
}
