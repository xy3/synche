package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/jobs"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
	"log"
	"os"
	"path"
	"sync"
)

var uploadCmd = &cobra.Command{
	Use:   "upload [file path]",
	Short: "Uploads a specified file to the server",
	Long:  `Uploads a specified local file to the server using chunked uploading`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uploadJob(args)
	},
}

var (
	fileName = ""
)

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.PersistentFlags().StringVarP(&fileName, "name", "n", "", "store the file " +
		"on the server with this name instead")
}

func uploadJob(args []string) {
	filePath := args[0]
	info, err := os.Stat(filePath)
	if err != nil {
		panic(err) // TODO
	}
	if info.IsDir() {
		log.Fatalf("The path specified is a directory: '%s'", filePath)
	}

	// Get the file name from the path if the --name flag is not set
	if len(fileName) == 0 {
		fileName = path.Base(filePath)
	}

	chunks, err := jobs.Split(filePath, viper.GetString("ChunkDir"))
	if err != nil {
		panic(err)
	}

	// Send a new file upload request to the server
	newUploadParams, err := upload.NewFileUploadParams(filePath, fileName, int64(len(chunks)))
	if err != nil {
		panic(err) // TODO
	}
	requestAccepted, err := upload.SendNewFileUploadRequest(newUploadParams)
	if err != nil {
		panic(err) // TODO
	}

	var wg sync.WaitGroup

	for chunkNum, chunk := range chunks {
		wg.Add(1)
		params, _ := upload.NewChunkUploadParams(chunk, "chunk_data", requestAccepted.UploadRequestID, int64(chunkNum))
		go upload.Chunk(&wg, params)
	}

	fmt.Println("Main: waiting for workers to finish")
	wg.Wait()
	fmt.Println("Main: Completed.")
}
