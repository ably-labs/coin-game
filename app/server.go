package app

import (
	"log"
	"net/http"
)

// RunServer - run server
func RunServer() (err error) {
	log.Println("Starting HTTP server on port 8000")
	server := NewRouter()
	a := server.InitializeRoutes()
	return http.ListenAndServe(":8000", a)
}
