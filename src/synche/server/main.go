package main

import (
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Received: "+r.FormValue("hash"))

	receivedFile, _, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer receivedFile.Close()

	newFile, err := os.Create("../data/received/" + r.FormValue("hash"))
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(newFile, receivedFile)
	if err != nil {
		panic(err)
	}
}


func main() {
	http.HandleFunc("/upload", uploadHandler)

	port := "5050"
    println("Listening on port " + port)
	http.ListenAndServe(":"+port, nil)

}
