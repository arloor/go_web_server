package main

import (
	"go_web_server/internal/server"
)

func main() {
	server.Serve()
	select {}
}
