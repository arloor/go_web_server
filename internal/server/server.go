package server

import (
	"go_web_server/internal/config"
	"net/http"
)

func Serve() error {
	http.HandleFunc("/ip", writeIp)
	http.HandleFunc("/", fileHandlerFunc())

	instance := config.Instance
	handler := MineHandler{}
	if !instance.UseTls {
		return http.ListenAndServe(instance.Addr, handler)
	} else {
		return http.ListenAndServeTLS(instance.Addr, instance.Cert, instance.PrivKey, handler)
	}
}
