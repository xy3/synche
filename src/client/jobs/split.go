package jobs

import (
	"encoding/hex"
	"fmt"
	"github.com/kalafut/imohash"
	"io/ioutil"
	"math"
	"os"
)

func Split(filePath, chunkDir string) ([]string, error) {

	file, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	fileSize := fileInfo.Size()

	const fileChunk = 1 * (1 << 20) // 1MB, change this to your requirement

	// calculate total number of parts the file will be chunked into

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	chunks := make([]string, totalPartsNum)

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		chunkSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		chunkBuffer := make([]byte, chunkSize)

		file.Read(chunkBuffer)

		// write to disk
		hash := imohash.Sum(chunkBuffer)
		fileName := chunkDir + "/" + hex.EncodeToString(hash[:])
		_, err := os.Create(fileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// write/save buffer to disk
		ioutil.WriteFile(fileName, chunkBuffer, os.ModeAppend)

		chunks[i] = fileName

		fmt.Println("Split to : ", fileName)
	}

	return chunks, err
}