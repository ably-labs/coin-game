package main

import (
	"log"
	server "stocker/app"
)

func main() {
	err := server.RunServer()
	if err != nil {
		log.Fatal("Server error")
	}

}
