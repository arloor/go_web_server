package server

import (
	"go_web_server/internal/config"
	"net/http"
	"time"
	"crypto/tls"
)

func Serve() error {
	http.HandleFunc("/ip", writeIp)
	http.HandleFunc("/", fileHandlerFunc())

	instance := config.Instance
	handler := MineHandler{}
	srv := &http.Server{
		Addr:              instance.Addr,
		Handler:           handler,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second, // Set idle timeout
		TLSConfig: &tls.Config{
			GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
                                               // Always get latest localhost.crt and localhost.key 
                                               // ex: keeping certificates file somewhere in global location where created certificates updated and this closure function can refer that
				cert, err := tls.LoadX509KeyPair(instance.Cert, instance.PrivKey)
				if err != nil {
					return nil, err
				}
				return &cert, nil
			},
		},
	}
	if !instance.UseTls {
		return srv.ListenAndServe()
	} else {
		return srv.ListenAndServeTLS("", "")
	}
}
