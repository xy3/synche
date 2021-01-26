package upload

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/kalafut/imohash"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func NewFileUploadRequest(url, filePath, fileParamName string) (*http.Request, error) {
	fileHash, err := imohash.SumFile(filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fileParamName, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)

	fileHashHex := hex.EncodeToString(fileHash[:])
	fmt.Println("Client-side hash:" + fileHashHex)

	err = writer.WriteField("hash", fileHashHex)
	if err != nil {
		panic(err)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, err
}



func SendUploadRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return response, err
	}

	return response, nil
}

func SendAndShowResult(request *http.Request) ServerResponse {
	response, _ := SendUploadRequest(request)

	body := &bytes.Buffer{}
	_, err := body.ReadFrom(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	serverResponse := ServerResponse{response.Status, response.Header, fmt.Sprint(body)}
	return serverResponse
}