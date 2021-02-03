package main

import (
	"fmt"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/testing"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/jobs"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
	"log"
	"os"
	"sync"
)

const (
	TestFilePath = "../data/test.mp4"
	ChunkDir     = "../data/chunks"
)

func check() {
	resp, err := apiclient.Default.Testing.CheckGet(testing.NewCheckGetParams())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", resp.Payload)
}

//// sample usage
func main() {
	// This function simply hits the testing endpoint of the server to check if it can connect to it
	check()

	// Create "chunks" directory if it doesn't exist
	_ = os.MkdirAll(ChunkDir, os.ModePerm)

	chunks, err := jobs.Split(TestFilePath, ChunkDir)
	if err != nil {
		panic(err)
	}

	// Send a new file upload request to the server
	newUploadParams, _ := upload.NewFileUploadParams(TestFilePath, "test.mp4", int64(len(chunks)))
	requestAccepted, err := upload.SendNewFileUploadRequest(newUploadParams)
	if err != nil {
		panic(err)
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
