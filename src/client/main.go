package main

import (
	"fmt"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/upload"
	"os"
	"sync"
)

func uploadChunk(wg *sync.WaitGroup, url string, chunkPath string) {
	defer wg.Done()
	request, err := upload.NewFileUploadRequest(url, chunkPath, "file")
	if err != nil {
		panic(err)
	}

	response := upload.SendAndShowResult(request)
	fmt.Println(response.Body)
}

// sample usage
func main() {
	url := "http://127.0.0.1:5050/upload"
	filePath := "../data/test.mp4"

	// Create "data" and "chunks" directory if they don't exist
	// os.ModePerm is equivalent to 0777
	if _, err := os.Stat("../data"); os.IsNotExist(err) {
		os.Mkdir("../data", os.ModePerm)
	}

	if _, err := os.Stat("../data/chunks"); os.IsNotExist(err) {
		os.Mkdir("../data/chunks", os.ModePerm)
	}
	
	chunks, err := files.Split(filePath, "../data/chunks")
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, chunk := range chunks {
		wg.Add(1)
		go uploadChunk(&wg, url, chunk)
	}

	fmt.Println("Main: waiting for workers to finish")
	wg.Wait()
	fmt.Println("Main: Completed.")
}
