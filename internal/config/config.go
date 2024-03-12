package config

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

type stringArray []string

func (i *stringArray) String() string {
	return fmt.Sprint(*i)
}

// Set 方法是flag.Value接口, 设置flag Value的方法.
// 通过多个flag指定的值， 所以我们追加到最终的数组上.
func (i *stringArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type Config struct {
	Addrs     stringArray       `yaml:"addrs"`
	UseTls    bool              `yaml:"tls"`
	Cert      string            `yaml:"cert"`
	PrivKey   string            `yaml:"key"`
	LogPath   string            `yaml:"log"`
	WebPath   string            `yaml:"content"`
	Users     stringArray       `yaml:"users"`
	BasicAuth map[string]string `yaml:"auth"`
	Refer     string            `yaml:"refer"`
}

var GlobalConfig Config

func init() {
	flag.Var(&GlobalConfig.Addrs, "addr", "监听地址，例如 :7788 。支持多个地址。默认监听 :7788")
	flag.BoolVar(&GlobalConfig.UseTls, "tls", false, "是否使用tls")
	flag.StringVar(&GlobalConfig.Cert, "cert", "cert.pem", "tls证书")
	flag.StringVar(&GlobalConfig.PrivKey, "key", "privkey.pem", "tls私钥")
	flag.StringVar(&GlobalConfig.LogPath, "log", "/tmp/proxy.log", "日志文件路径")
	flag.StringVar(&GlobalConfig.WebPath, "content", ".", "文件服务器目录")
	flag.Var(&GlobalConfig.Users, "user", "Basic认证的用户名密码，例如username:password")
	flag.StringVar(&GlobalConfig.Refer, "refer", "", "本站的referer特征")
	flag.Parse()
	if len(GlobalConfig.Addrs) == 0 {
		GlobalConfig.Addrs = append(GlobalConfig.Addrs, ":7788")
	}
	GlobalConfig.BasicAuth = make(map[string]string)
	for _, user := range GlobalConfig.Users {
		base64Encode := "Basic " + base64.StdEncoding.EncodeToString([]byte(user))
		GlobalConfig.BasicAuth[base64Encode] = strings.Split(user, ":")[0]
	}
	initLog()
	out, err := yaml.Marshal(GlobalConfig)
	if err != nil {
		log.Println("go web server config:", GlobalConfig)
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
	file := GlobalConfig.LogPath
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
