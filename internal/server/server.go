package server

import (
	"go_web_server/internal/config"
	"net/http"
	"time"
)

func Serve() error {
	http.HandleFunc("/ip", writeIp)
	http.HandleFunc("/", fileHandlerFunc())

	instance := config.Instance
	handler := MineHandler{}
	srv := &http.Server{
		Addr:         instance.Addr,
		Handler:      handler,
		WriteTimeout: 30 * time.Second, // Set idle timeout
	}
	if !instance.UseTls {
		return srv.ListenAndServe()
	} else {
		return srv.ListenAndServeTLS(instance.Cert, instance.PrivKey)
	}
}
