package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
)

func NewUploadCmd(fileUploader upload.FileUploader) *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:   "upload [file path]",
		Short: "Uploads a specified file to the server",
		Long:  `Uploads a specified local file to the server using chunked uploading`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			return fileUploader.Upload(filePath)
		},
	}
	uploadCmd.Flags().StringP("name", "n", "", "store the file on the server with this name instead")
	return uploadCmd
}

func init() {
	fileUploader := upload.NewFileUpload(*data.NewSplitter(data.DefaultChunkWriter), *upload.NewChunkUploader(apiclient.Default.Files.UploadChunk))
	uploadCmd := NewUploadCmd(fileUploader)
	rootCmd.AddCommand(uploadCmd)

	err := viper.BindPFlags(uploadCmd.Flags())
	if err != nil {
		log.Fatalf("Could not bind flags to viper config: %v", err)
	}
}
