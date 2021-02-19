package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
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
	splitter      data.Splitter
	chunkUploader upload.ChunkUploader
	fileUploader  upload.FileUploader
	chunksDir     string
}

func (u UploadJob) Run(filePath string) error {
	file, err := data.AppFS.Open(filePath)
	if err != nil {
		return nil
	}
	return u.fileUploader.Upload(file)
}

func NewUploadJob(
	chunkWriter data.ChunkWriter,
	chunkHashFunc data.ChunkHashFunc,
	fileHashFunc data.FileHashFunc,
	newUploadRequester upload.NewUploadRequester,
	chunksDir string,
	chunkMBs uint64,
) (
	*UploadJob,
) {
	splitter := data.NewSplitJob(chunkWriter, chunkHashFunc, chunksDir, chunkMBs)
	chunkUploader := upload.NewChunkUpload(chunkHashFunc)
	fileUploader := upload.NewFileUpload(splitter, chunkUploader, fileHashFunc, newUploadRequester)

	return &UploadJob{
		chunksDir:     chunksDir,
		splitter:      splitter,
		chunkUploader: chunkUploader,
		fileUploader:  fileUploader,
	}
}

func NewDefaultUploadJob() *UploadJob {
	return NewUploadJob(data.DefaultChunkWriter, data.DefaultChunkHashFunc, data.DefaultFileHashFunc, upload.DefaultNewUploadRequester, c.Config.Chunks.Dir, 1)
}

func init() {
	uploadCmd := NewUploadCmd(NewDefaultUploadJob())
	rootCmd.AddCommand(uploadCmd)

	err := viper.BindPFlags(uploadCmd.Flags())
	if err != nil {
		log.Fatalf("Could not bind flags to viper config: %v", err)
	}
}
