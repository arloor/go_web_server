package server

import (
	"go_web_server/internal/config"
	"net/http"
)

func Serve() {
	http.HandleFunc("/ip", writeIp)
	instance := config.Instance
	if !instance.UseTls {
		http.ListenAndServe(instance.Addr, nil)
	} else {
		http.ListenAndServeTLS(instance.Addr, instance.Cert, instance.PrivKey, nil)
	}
}
