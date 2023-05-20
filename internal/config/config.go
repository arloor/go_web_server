package config

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

type Config struct {
	Addr    string
	UseTls  bool
	Cert    string
	PrivKey string
	LogPath string
}

var Instance Config

const AddrEnv = "go_server_port"
const UseTls = "use_tls"
const CERT = "cert"
const KEY = "key"
const logPath = "logPath"

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
	logfile := os.Getenv(logPath)
	if logfile == "" {
		logfile = "/var/log/go_web_server.log"
	}
	Instance = Config{
		Addr:    addrEnv,
		UseTls:  useTls,
		Cert:    cert,
		PrivKey: key,
		LogPath: logfile,
	}
	initLog()
	log.Println("go web server config:", Instance)
}

func initLog() {
	file := Instance.LogPath
	rollingFile := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    50,
		MaxAge:     14,
		MaxBackups: 10,
		Compress:   false,
	}
	mw := io.MultiWriter(os.Stdout, rollingFile)
	log.SetOutput(mw)
	log.SetFlags(log.Lshortfile | log.Flags())
}
