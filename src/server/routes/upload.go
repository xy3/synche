package routes

import (
	"io"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Check that the required parameters are provided in the request
	io.WriteString(w, "Received: "+r.FormValue("hash"))

	receivedFile, _, err := r.FormFile("file")
	if err != nil {
		panic(err) // todo
	}
	defer receivedFile.Close()

	// TODO: manage the chunk data received in the request
	newFile, err := os.Create("../data/received/" + r.FormValue("hash"))
	if err != nil {
		panic(err) // todo
	}

	_, err = io.Copy(newFile, receivedFile)
	if err != nil {
		panic(err) // todo
	}
}