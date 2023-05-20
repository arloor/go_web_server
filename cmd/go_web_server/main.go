package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	portStr := os.Getenv("go_server_port")
	if portStr == "" {
		portStr = ":8080"
	}
	fmt.Fprintf(os.Stdout, "Server is running on port %s ...", portStr)
	http.HandleFunc("/ip", writeIp)
	http.ListenAndServe(portStr, nil)
	select {}
}

func writeIp(w http.ResponseWriter, r *http.Request) {
	remoteAddr := r.RemoteAddr
	index := strings.LastIndex(remoteAddr, ":")
	if index == -1 {
		fmt.Fprintf(w, "addr is not ip:port %s", remoteAddr)
	} else {
		ip := remoteAddr[:index]
		fmt.Fprint(w, ip)
	}

}
