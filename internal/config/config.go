package config

import (
	"flag"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

type Config struct {
	Addr      string `yaml:"addr"`
	UseTls    bool   `yaml:"tls"`
	Cert      string `yaml:"cert"`
	PrivKey   string `yaml:"key"`
	LogPath   string `yaml:"log"`
	WebPath   string `yaml:"content"`
	BasicAuth string `yaml:"auth"`
}

var Instance Config

func init() {
	flag.StringVar(&Instance.Addr, "addr", ":7777", "监听地址")
	flag.BoolVar(&Instance.UseTls, "tls", false, "是否使用tls")
	flag.StringVar(&Instance.Cert, "cert", "cert.pem", "tls证书")
	flag.StringVar(&Instance.PrivKey, "key", "privkey.pem", "tls私钥")
	flag.StringVar(&Instance.LogPath, "log", "/tmp/proxy.log", "日志文件路径")
	flag.StringVar(&Instance.WebPath, "content", "/data", "文件服务器目录")
	flag.StringVar(&Instance.BasicAuth, "auth", "", "Basic Auth Header")
	flag.Parse()
	initLog()
	out, err := yaml.Marshal(Instance)
	if err != nil {
		log.Println("go web server config:", Instance)
	} else {
		log.Printf("go web server config: \n%s", string(out))
	}

}

//const AddrEnv = "addr"
//const UseTls = "use_tls"
//const CERT = "cert"
//const KEY = "key"
//const logPath = "log_path"
//const webPath = "web_path"
//const constBasicAuth = "basic_auth"
//
//func init() {
//	useTls := os.Getenv(UseTls) == "true"
//	addrEnv := os.Getenv(AddrEnv)
//	if addrEnv == "" {
//		if useTls {
//			addrEnv = ":8443"
//		} else {
//			addrEnv = ":8080"
//		}
//	}
//	cert := os.Getenv(CERT)
//	if cert == "" {
//		cert = "cert.pem"
//	}
//	key := os.Getenv(KEY)
//	if key == "" {
//		key = "privkey.pem"
//	}
//	logfile := os.Getenv(logPath)
//	if logfile == "" {
//		logfile = "/var/log/go_web_server.log"
//	}
//	webContentPath := os.Getenv(webPath)
//	if webContentPath == "" {
//		webContentPath = "."
//	}
//	basicAuth := os.Getenv(constBasicAuth)
//	Instance = Config{
//		Addr:      addrEnv,
//		UseTls:    useTls,
//		Cert:      cert,
//		PrivKey:   key,
//		LogPath:   logfile,
//		WebPath:   webContentPath,
//		BasicAuth: basicAuth,
//	}
//	initLog()
//	log.Println("go web server config:", Instance)
//}

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
