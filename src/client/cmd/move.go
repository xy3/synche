package cmd

import (
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

// printFile Logs file details in a readable format
func printFile(file *models.File) {
	log.Info("--- File moved successfully. Updated file information: ---")
	log.WithFields(log.Fields{
		"Name":        file.Name,
		"ID":          file.ID,
		"DirectoryID": file.DirectoryID,
		"FileSize":    file.Size,
	}).Infof(color.GreenString(file.Name))
}

func moveFileByID(fileID uint64, fileUpdate *models.FileUpdate) (*models.File, error) {
	params := &files.UpdateFileByIDParams{
		FileID:     fileID,
		FileUpdate: fileUpdate,
	}

	resp, err := apiclient.Client.Files.UpdateFileByID(params, apiclient.ClientAuth)
	if err != nil {
		return nil, err
	}
	return resp.GetPayload(), nil
}

func moveFileByPath(filePath string, fileUpdate *models.FileUpdate) (*models.File, error) {
	params := &files.UpdateFileByPathParams{
		FilePath:   filePath,
		FileUpdate: fileUpdate,
	}

	resp, err := apiclient.Client.Files.UpdateFileByPath(params, apiclient.ClientAuth)
	if err != nil {
		return nil, err
	}
	return resp.GetPayload(), nil
}

var (
	moveCurrentFilePath string
	moveNewFilePath     string
	moveFileID          uint64
	moveDirID           uint64
	moveNewFileName     string
)

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move a file",
	Aliases: []string{"mv", "move"},
	Long: `Move a file from one specified location to another using the full 
path to the current location or file ID, and the full path to the new 
location or the directory ID`,
	PreRun: authenticateUserPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err  error
			file *models.File
		)

		fileUpdate := &models.FileUpdate{
			NewDirectoryID: moveDirID,
			NewFileName:    moveNewFileName,
			NewFilePath:    moveNewFilePath,
		}

		if moveFileID != 0 {
			file, err = moveFileByID(moveFileID, fileUpdate)
		} else if moveCurrentFilePath != "" {
			file, err = moveFileByPath(moveCurrentFilePath, fileUpdate)
		} else {
			_ = cmd.Usage()
		}

		if err != nil {
			log.WithError(err).Fatal("Failed to move file")
		}
		printFile(file)
	},
}

// TODO: Fix bug so that arguments can default to file paths
func init() {
	moveCmd.Flags().StringVarP(&moveCurrentFilePath, "file-path", "f", "", "the file to move")
	moveCmd.Flags().Uint64VarP(&moveFileID, "file-id", "i", 0, "the ID of the file to move")
	moveCmd.Flags().StringVarP(&moveNewFilePath, "new-file-path", "o", "", "the new path to move to")
	moveCmd.Flags().Uint64VarP(&moveDirID, "dir-id", "d", 0, "the directory ID to move to")
	rootCmd.AddCommand(moveCmd)
}
