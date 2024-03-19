package server

import "net/http"

type MineHandler struct {
}

func (MineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.Method == http.MethodConnect { //connect处理
		connect(w, r)
	} else {
		http.DefaultServeMux.ServeHTTP(w, r)
	}

}
