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
	"strconv"
)

func printFile(file *models.File) {
	log.Infof("--- File moved successfully. Updated file information: ---")
	log.WithFields(log.Fields{
		"Name":			*file.Name,
		"ID":           *file.ID,
		"DirectoryID":  *file.StorageDirectoryID,
		"FileSize":     *file.Size}).Infof(color.GreenString(*file.Name))
}

func moveJob(params *files.MoveFileParams) *files.MoveFileOK {
	resp, err := apiclient.Client.Files.MoveFile(params, apiclient.ClientAuth)
	if err != nil {
		log.WithError(err).Fatal("Failed to move file")
	}
	return resp
}

func moveToDirectoryID(fileID uint64) *models.File {
	resp := moveJob(files.NewMoveFileParams().WithFileID(fileID).WithDirectoryID(&moveDirID))
	log.Infof("Directory name: %v\n", &moveDirName)
	return resp.GetPayload()
}

func moveToDirectoryName(fileID uint64) *models.File {
	resp := moveJob(files.NewMoveFileParams().WithFileID(fileID).WithDirectoryName(&moveDirName))
	log.Infof("Directory name: %v\n", &moveDirName)
	return resp.GetPayload()
}

var moveDirName string
var moveDirID uint64
var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "move FileID DirectoryID",
	Long:  `<Move a file from one specified location to another using the file ID and directory IDs>`,
	Run: func(cmd *cobra.Command, args []string) {
		fileID, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Error("Invalid file ID")
		}

		err = apiclient.Authenticator(filepath.Join(config.SyncheDir, "token.json"))
		if err != nil {
			log.Fatal("Failed to authenticate the client")
		}

		var updatedFile *models.File
		if moveDirName != "" {
			updatedFile = moveToDirectoryName(fileID)
		} else {
			updatedFile = moveToDirectoryID(fileID)
		}

		printFile(updatedFile)
	},
}

func init() {
	moveCmd.Flags().StringVarP(&moveDirName, "dir-name", "n", "", "Specify the name of a directory to move a file to it")
	moveCmd.Flags().Uint64VarP(&moveDirID, "dir-id", "i", 0, "Specify an ID of a directory to move a file to it")
	rootCmd.AddCommand(moveCmd)
}