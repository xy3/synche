package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
	"path"
)

func NewUploadCmd(uploader Uploader) *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:   "upload [file path]",
		Short: "Uploads a specified file to the server",
		Long:  `Uploads a specified local file to the server using chunked uploading`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filePath := args[0]
			err := uploader.Run(filePath)
			if err != nil {
				log.Fatalf("Failed to upload the file: %v", err)
			}
		},
	}
	uploadCmd.Flags().StringP("name", "n", "", "store the file on the server with this name instead")
	return uploadCmd
}

//go:generate mockery --name=Uploader --case underscore
type Uploader interface {
	Run(filePath string) error
}

type UploadJob struct {
	chunkUploader upload.ChunkUploader
	fileUploader  upload.FileUploader
	fileHashFunc  data.FileHashFunc
}

func (u UploadJob) Run(filePath string) error {
	file, err := files.AppFS.Open(filePath)
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	hash, err := u.fileHashFunc(filePath)
	if err != nil {
		return err
	}
	splitFile := data.NewSplitFile(stat.Size(), c.Config.Chunks.Size, filePath, path.Base(filePath), hash, file)
	return u.fileUploader.Upload(splitFile)
}

func NewUploadJob(newUploadRequester upload.NewUploadRequester, fileHashFunc data.FileHashFunc) *UploadJob {
	chunkUploader := new(upload.ChunkUpload)
	fileUploader := upload.NewFileUpload(chunkUploader, newUploadRequester)

	return &UploadJob{
		chunkUploader: chunkUploader,
		fileUploader:  fileUploader,
		fileHashFunc:  fileHashFunc,
	}
}

func NewDefaultUploadJob() *UploadJob {
	return NewUploadJob(upload.DefaultNewUploadRequester, data.DefaultFileHashFunc)
}

func init() {
	uploadCmd := NewUploadCmd(NewDefaultUploadJob())
	rootCmd.AddCommand(uploadCmd)

	err := viper.BindPFlags(uploadCmd.Flags())
	if err != nil {
		log.Fatalf("Could not bind flags to viper config: %v", err)
	}
}
