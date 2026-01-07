package server

import (
	"crypto/tls"
	"go_web_server/internal/config"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var sslCert *tls.Certificate = nil
var sslLastCertUpdateTime time.Time = time.Now()

const sslCertUpdateInterval = 5 * time.Hour

var (
	ReqCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "req_from_out_total",
		Help: "Number of HTTP requests received",
	}, []string{"referer", "path"})
	ProxyTraffic = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "proxy_traffic_total",
		Help: "num proxy_traffic",
	}, []string{"client", "target", "username"})
)

func Serve() error {
	http.HandleFunc("/ip", writeIP)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", fileHandlerFunc())

	errors := make(chan error)

	globalConfig := config.GlobalConfig
	handler := MineHandler{}
	for _, addr := range globalConfig.Addrs {
		srv := &http.Server{
			Addr:              addr,
			Handler:           handler,
			IdleTimeout:       31 * time.Second,
			ReadHeaderTimeout: 31 * time.Second,
			ReadTimeout:       31 * time.Second,
			WriteTimeout:      31 * time.Second, // Set idle timeout
			TLSConfig: &tls.Config{
				GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
					return loadNewCertIfNeed(globalConfig.Cert, globalConfig.PrivKey)
				},
			},
		}
		go func() {
			if !globalConfig.UseTlS {
				errors <- srv.ListenAndServe()
			} else {
				errors <- srv.ListenAndServeTLS("", "")
			}
		}()
	}
	return <-errors
}

func loadNewCertIfNeed(certFile, privkey string) (*tls.Certificate, error) {
	now := time.Now()
	if sslCert == nil || now.Sub(sslLastCertUpdateTime) > sslCertUpdateInterval {
		cert, err := tls.LoadX509KeyPair(certFile, privkey)
		if err != nil {
			log.Println("Error loading certificate", err)
			if sslCert != nil {
				return sslCert, nil
			}
			return nil, err
		} else {
			log.Println("Loaded certificate", certFile, privkey)
		}
		sslCert = &cert
		sslLastCertUpdateTime = now
		return &cert, nil
	} else {
		return sslCert, nil
	}
}
