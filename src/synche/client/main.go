package main

import (
	"client/files"
	"client/upload"
	"fmt"
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