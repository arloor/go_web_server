package server

import (
	"go_web_server/internal/config"
	"net/http"
)

func serve() {
	config := config.Instance
	http.HandleFunc("/ip", writeIp)
	if !config.UseTls {
		http.ListenAndServe(config.Addr, nil)
	}else {
		http.ListenAndServe(config.Addr,config.CERT,config., nil)
	}
}


