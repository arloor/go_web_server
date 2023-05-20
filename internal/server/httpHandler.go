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

func fileHandler() http.Handler {
	fs := http.FileServer(http.Dir("."))

	// 设定路由，所有的请求都交给fs去处理
	return http.StripPrefix("/", fs)
}
