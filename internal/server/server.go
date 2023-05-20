package server

import (
	"go_web_server/internal/config"
	"net/http"
)

func Serve() error {
	http.HandleFunc("/ip", writeIp)
	http.HandleFunc("/", fileHandler().ServeHTTP)
	instance := config.Instance
	if !instance.UseTls {
		return http.ListenAndServe(instance.Addr, nil)
	} else {
		return http.ListenAndServeTLS(instance.Addr, instance.Cert, instance.PrivKey, nil)
	}
	return nil
}
