package server

import (
	"errors"
	"go_web_server/internal/config"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
)

func dialContextCheckACL(network, hostPort string) (net.Conn, error) {
	// This is net.Dial's default behavior: if the host resolves to multiple IP addresses,
	// Dial will try each IP address in order until one succeeds
	return net.Dial(network, hostPort)
}

func connect(w http.ResponseWriter, r *http.Request) {
	clientAddr := strings.Split(r.RemoteAddr, ":")[0]
	var username string
	if len(config.GlobalConfig.BasicAuth) != 0 {
		var ok bool
		username, ok = config.GlobalConfig.BasicAuth[r.Header.Get("proxy-authorization")]
		if !ok {
			log.Println("wrong basic auth from", clientAddr)
			http.Error(w, "InternalServerError", http.StatusInternalServerError)
			return
		}
	}
	if r.ProtoMajor == 2 {
		if len(r.URL.Scheme) > 0 || len(r.URL.Path) > 0 {
			http.Error(w, "CONNECT request has :scheme or/and :path pseudo-header fields", http.StatusBadRequest)
			return
		}
	}

	hostPort := r.URL.Host
	if hostPort == "" {
		hostPort = r.Host
	}
	targetConn, err := dialContextCheckACL("tcp", hostPort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if targetConn == nil {
		// safest to check both error and targetConn afterwards, in case fp.dial (potentially unstable
		// from x/net/proxy) misbehaves and returns both nil or both non-nil
		http.Error(w, "hostname "+r.URL.Hostname()+" is not allowed", http.StatusBadRequest)
		return
	}
	defer targetConn.Close()

	switch r.ProtoMajor {
	case 1: // http1: hijack the whole flow
		_, err := serveHijack(w, targetConn, clientAddr, hostPort, username)
		if err != nil {
			log.Println(err, r.RemoteAddr)
		}
		return
	case 2: // http2: keep reading from "request" and writing into same response
		defer r.Body.Close()
		wFlusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "ResponseWriter doesn't implement Flusher()", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		randN := rand.Intn(150)
		for i := 0; i < randN; i++ {
			w.Header().Add("Server", "go_web_server")
		}
		wFlusher.Flush()
		err := dualStream(targetConn, r.Body, w, clientAddr, hostPort, username)
		if err != nil {
			log.Println(err, r.RemoteAddr)
		}
		return
	default:
		http.Error(w, "Wrong http version", http.StatusBadRequest)
	}
}

// Hijacks the connection from ResponseWriter, writes the response and proxies data between targetConn
// and hijacked connection.
func serveHijack(w http.ResponseWriter, targetConn net.Conn, clientAddr string, hostPort string, username string) (int, error) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		return http.StatusInternalServerError, errors.New("ResponseWriter does not implement Hijacker")
	}
	clientConn, bufReader, err := hijacker.Hijack()
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to hijack: " + err.Error())
	}
	defer clientConn.Close()
	// bufReader may contain unprocessed buffered data from the client.
	if bufReader != nil {
		// snippet borrowed from `proxy` plugin
		if n := bufReader.Reader.Buffered(); n > 0 {
			rbuf, peekErr := bufReader.Reader.Peek(n)
			if peekErr != nil {
				return http.StatusBadGateway, peekErr
			}
			if _, writeErr := targetConn.Write(rbuf); writeErr != nil {
				return http.StatusBadGateway, writeErr
			}
		}
	}
	// Since we hijacked the connection, we lost the ability to write and flush headers via w.
	// Let's handcraft the response and send it manually.
	res := &http.Response{StatusCode: http.StatusOK,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
	}
	for i := 0; i < rand.Intn(150); i++ {
		res.Header.Add("Server", "go_web_server")
	}

	err = res.Write(clientConn)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to send response to client: " + err.Error())
	}

	return 0, dualStream(targetConn, clientConn, clientConn, clientAddr, hostPort, username)
}

var bufferPool = &sync.Pool{New: func() interface{} {
	return make([]byte, 32*1024)
}}

func dualStream(target net.Conn, clientReader io.ReadCloser, clientWriter io.Writer, clientAddr string, hostPort string, username string) error {
	stream := func(w io.Writer, r io.Reader) error {
		// copy bytes from r to w
		buf, ok := bufferPool.Get().([]byte)
		if !ok {
			return errors.New("failed to get buffer from pool")
		}
		defer bufferPool.Put(buf) //nolint:staticcheck
		buf = buf[0:cap(buf)]
		nw, _err := flushingIoCopy(w, r, buf)
		ProxyTraffic.WithLabelValues(clientAddr, hostPort, username).Add(float64(nw))
		if closeWriter, ok := w.(interface {
			CloseWrite() error
		}); ok {
			if err := closeWriter.CloseWrite(); err != nil {
				return err
			}
		}
		return _err
	}

	go func() {
		if err := stream(target, clientReader); err != nil {
			log.Println("Error in stream:", err)
		}
	}()
	return stream(clientWriter, target)
}

func flushingIoCopy(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	flusher, ok := dst.(http.Flusher)
	if !ok {
		return io.CopyBuffer(dst, src, buf)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			flusher.Flush()
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return
}
