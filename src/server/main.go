package main

import (
	"net/http"
)

func init() {

}

func main() {
	port := "5050"
    println("Listening on port " + port)
	http.ListenAndServe(":"+port, nil)
}
