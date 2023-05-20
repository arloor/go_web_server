package server

import (
	"fmt"
	"net/http"
	"strings"
)

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
