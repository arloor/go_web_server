## static file serve AND Http TUNNEL proxy

### Usage

```bash
$ curl -LO https://github.com/arloor/go_web_server/releases/download/v1.0.0/go_web_server # v1.0.0 is auto updated by github action
$ chmod +x go_web_server                                                                 
$ ./go_web_server 
2023/12/24 15:46:27 config.go:38: go web server config: 
addr: :8080
tls: false
cert: cert.pem
key: privkey.pem
log: /tmp/proxy.log
content: /data
auth: ""
2023/12/24 15:46:35 httpHandlerFunc.go:47: xxxxxxxxx:55054     GET / HTTP/1.1
2023/12/24 15:46:35 httpHandlerFunc.go:47: xxxxxxxxx:55054     GET /favicon.ico HTTP/1.1
2023/12/24 15:46:36 httpHandlerFunc.go:47: xxxxxxxxx:55054     GET /rust_http_proxy/ HTTP/1.1
.....
2023/12/24 15:47:38 httpHandlerFunc.go:45: xxxxxxxxx:56902 CONNECT baidu.com:443 HTTP/1.1
2023/12/24 15:47:46 httpHandlerFunc.go:45: xxxxxxxxx:38208 CONNECT www.baidu.com:443 HTTP/1.1
```

### static files dowload

![Alt text](image.png)

### Http TUNNEL proxy

```bash
curl https://www.baidu.com -p -x http:/your_addr:8080 -v
```
