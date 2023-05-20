package main

import (
	"go_web_server/internal/server"
	"log"
)

func main() {
	err := server.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
