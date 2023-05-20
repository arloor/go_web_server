package server

import (
	"fmt"
	"go_web_server/internal/config"
	"net/http"
)

func Serve() {
	instance := config.Instance
	http.HandleFunc("/ip", writeIp)
	if !instance.UseTls {
		fmt.Printf("listen on http")
		http.ListenAndServe(instance.Addr, nil)
	} else {
		http.ListenAndServeTLS(instance.Addr, instance.Cert, instance.PrivKey, nil)
	}
}
