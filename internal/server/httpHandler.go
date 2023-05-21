package server

import (
	"fmt"
	"go_web_server/internal/config"
	"log"
	"net/http"
	"strings"
)

func writeIp(w http.ResponseWriter, r *http.Request) {
	log.Println("request for", r.URL.Path, "from", r.RemoteAddr)
	remoteAddr := r.RemoteAddr
	index := strings.LastIndex(remoteAddr, ":")
	if index == -1 {
		fmt.Fprintf(w, "addr is not ip:port %s", remoteAddr)
	} else {
		ip := remoteAddr[:index]
		fmt.Fprint(w, ip)
	}

}

func fileHandlerFunc() http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		if containsDotDot(r.URL.Path) {
			// Too many programs use r.URL.Path to construct the argument to
			// serveFile. Reject the request under the assumption that happened
			// here and ".." may not be wanted.
			// Note that name might not contain "..", for example if code (still
			// incorrectly) used filepath.Join(myDir, r.URL.Path).
			log.Println(r.URL.Path, "is invalid from", r.RemoteAddr)
			http.Error(w, "invalid URL path", http.StatusBadRequest)
			return
		}
		fs := http.FileServer(http.Dir(config.Instance.WebPath))
		http.StripPrefix("/", fs).ServeHTTP(w, r)
	})
}

func logRequest(r *http.Request) {
	log.Println("request for", r.URL.Path, "from", r.RemoteAddr)
}

func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }
