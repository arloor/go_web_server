package server

import (
	"crypto/tls"
	"go_web_server/internal/config"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ssl_cert *tls.Certificate = nil
var ssl_last_cert_update time.Time = time.Now()

const ssl_cert_update_interval = 5 * time.Hour

var (
	HttpRequst = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "req_from_out_total",
		Help: "Number of HTTP requests received",
	}, []string{"referer", "path"})
	ProxyTraffic = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "proxy_traffic_total",
		Help: "num proxy_traffic",
	}, []string{"client", "target"})
)

func Serve() error {
	http.HandleFunc("/ip", writeIp)
	http.Handle("/metrics", promhttp.Handler())
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
				now := time.Now()
				if ssl_cert == nil || now.Sub(ssl_last_cert_update) > ssl_cert_update_interval {
					cert, err := tls.LoadX509KeyPair(instance.Cert, instance.PrivKey)
					if err != nil {
						log.Println("Error loading certificate", err)
						if ssl_cert != nil {
							return ssl_cert, nil
						}
						return nil, err
					} else {
						log.Println("Loaded certificate", instance.Cert, instance.PrivKey)
					}
					ssl_cert = &cert
					ssl_last_cert_update = now
					return &cert, nil
				} else {
					return ssl_cert, nil
				}
			},
		},
	}
	if !instance.UseTls {
		return srv.ListenAndServe()
	} else {
		return srv.ListenAndServeTLS("", "")
	}
}
