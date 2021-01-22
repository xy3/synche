package main

import (
	"net/http"
	"server/routes"
)

func init() {

}

func main() {
	routes.SetupRoutes()
	port := "5050"
    println("Listening on port " + port)
	http.ListenAndServe(":"+port, nil)
}
