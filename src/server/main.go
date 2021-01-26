package main

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/routes"
	"net/http"
)

func init() {

}

func main() {
	routes.SetupRoutes()
	port := "5050"
    println("Listening on port " + port)
	http.ListenAndServe(":"+port, nil)
}
