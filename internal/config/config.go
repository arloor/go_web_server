package config

import (
	"os"
)

type Config struct {
	Addr    string
	UseTls  bool
	Cert    string
	PrivKey string
}

var Instance Config

const AddrEnv = "go_server_port"
const UseTls = "use_tls"
const CERT = "cert"
const KEY = "key"

func init() {
	useTls := os.Getenv(UseTls) == "true"
	addrEnv := os.Getenv(AddrEnv)
	if addrEnv == "" {
		if useTls {
			addrEnv = ":8443"
		} else {
			addrEnv = ":8080"
		}
	}
	cert := os.Getenv(CERT)
	if cert == "" {
		cert = "cert.pem"
	}
	key := os.Getenv(KEY)
	if key == "" {
		key = "privkey.pem"
	}
	Instance = Config{
		Addr:    addrEnv,
		UseTls:  useTls,
		Cert:    cert,
		PrivKey: key,
	}
}
