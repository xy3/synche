package routes

import "net/http"

func SetupRoutes() {
	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/list", ListHandler)
	http.HandleFunc("/delete", DeleteHandler)
}
