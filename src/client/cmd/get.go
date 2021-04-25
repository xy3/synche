package cmd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	files2 "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/transfer"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"os"
	"path/filepath"
)

var fileID uint64
var filename string
var output string

func NewGetCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Download a file from a Synche server",
		Long: `Specify a file ID to download. The file will be downloaded to the current directory by default. Examples:
  synche get -i 2
  synche get -i 2 -o downloaded_file.jpg
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := apiclient.Authenticator(filepath.Join(config.SyncheDir, "token.json"))
			if err != nil {
				log.WithError(err).Fatal("Failed to authenticate the client")
			}
			log.Info("Downloading...")
			err = getFileJob(cmd)
			if err != nil {
				log.WithError(err).Fatal("download failed")
			}
		},
	}

	getCmd.Flags().Uint64VarP(&fileID, "id", "i", 0, "ID of the file to download")
	// getCmd.Flags().StringVarP(&filename, "name", "n", "", "name of the file to download")
	getCmd.Flags().StringVarP(&output, "output", "o", "", "download location. either a full file path or directory")
	return getCmd
}

func getFileJob(cmd *cobra.Command) error {
	if cmd.Flags().NFlag() == 0 {
		return cmd.Help()
	}

	fileInfo, err := apiclient.Client.Files.GetFileInfo(
		&files2.GetFileInfoParams{
			FileID:  fileID,
			Context: context.Background(),
		},
		apiclient.ClientAuth,
	)

	if err != nil {
		return err
	}
	if fileInfo.Payload.ID == nil {
		return fmt.Errorf("no file found for ID: %d", fileID)
	}

	var outputFile = *fileInfo.Payload.Name
	if output != "" {
		outputFile = output
	}
	fileWriter, err := files.Afs.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer fileWriter.Close()

	params := &transfer.DownloadFileParams{
		FileID:  int64(fileID),
		Context: context.Background(),
	}

	_, err = apiclient.Client.Transfer.DownloadFile(params, apiclient.ClientAuth, fileWriter)
	if err != nil {
		return err
	}

	log.Infof("Finished downloading the file to: %s", outputFile)
	return nil
}

func init() {
	rootCmd.AddCommand(NewGetCmd())
}
