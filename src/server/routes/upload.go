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

	// Create "received" and "data" directory if they don't exist
	// When the connection request feature is implemented, this should be done there instead
	// os.ModePerm is equivalent to 0777
	if _, err := os.Stat("../data"); os.IsNotExist(err) {
		os.Mkdir("../data", os.ModePerm)
	}

	if _, err := os.Stat("../data/received"); os.IsNotExist(err) {
		os.Mkdir("../data/received", os.ModePerm)
	}

	newFile, err := os.Create("../data/received/" + r.FormValue("hash"))
	if err != nil {
		panic(err) // todo
	}

	_, err = io.Copy(newFile, receivedFile)
	if err != nil {
		panic(err) // todo
	}
}