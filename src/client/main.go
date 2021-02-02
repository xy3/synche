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

func uploadChunk(wg *sync.WaitGroup, chunkPath string) {
	defer wg.Done()

	params, _ := upload.NewChunkUploadParams(chunkPath, "chunk_data")
	resp, err := apiclient.Default.Files.UploadFile(params)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", resp.Payload)
}

func check() {
	resp, err := apiclient.Default.Testing.CheckGet(testing.NewCheckGetParams())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", resp.Payload)
}

//// sample usage
func main() {
	// This function simple hits the testing endpoint of the server to check if it can connect to it
	check()

	filePath := "../data/test.mp4"

	// Create "data" and "chunks" directory if they don't exist
	// os.ModePerm is equivalent to 0777
	if _, err := os.Stat("../data"); os.IsNotExist(err) {
		os.Mkdir("../data", os.ModePerm)
	}

	if _, err := os.Stat("../data/chunks"); os.IsNotExist(err) {
		os.Mkdir("../data/chunks", os.ModePerm)
	}

	chunks, err := jobs.Split(filePath, "../data/chunks")

	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, chunk := range chunks {
		wg.Add(1)
		go uploadChunk(&wg, chunk)
	}

	fmt.Println("Main: waiting for workers to finish")
	wg.Wait()
	fmt.Println("Main: Completed.")
}
