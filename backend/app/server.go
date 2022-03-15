package app

import (
	"fmt"
	"log"
	"net/http"
	"stocker/config"
)

// RunServer - run server
func RunServer() (err error) {
	port := config.EnvVariable("PORT")
	if port == "" {
		port = "9000"
	}
	log.Printf("Starting HTTP server on port %s", port)
	server := NewRouter()
	a := server.InitializeRoutes()
	fmt.Println(port, "++")

	return http.ListenAndServe(fmt.Sprintf(":%s", port), a)
}
