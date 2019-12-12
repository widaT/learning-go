# go http2 开发

## 开发使用的https证书

https证书对开发人员不是很友好，安装开发环境的证书都相对麻烦了点。这边介绍一个小工具[mkcert](https://github.com/FiloSottile/mkcert)。

### 安装 mkcert
本文以ubuntu和debain系统为例：

- 安装依赖 `sudo apt install libnss3-tools`
- 下载[mkcert](https://github.com/FiloSottile/mkcert/releases/download/v1.4.1/mkcert-v1.4.1-linux-amd64)
- 生产证书`./mkcert -key-file key.pem -cert-file cert.pem example.com *.example.com localhost 127.0.0.1 ::1`
- 添加本地信任 `mkcert -install`

## golang http2

go 官方的http2包在[拓展网络库中](golang.org/x/net/http2),使用起来标准库中的http没啥区别。需要主要的是http2的包同时支持http2和http1当发现client端不支持http2的时候会使用http1。

我们新建一个项目http2test，同时把上面生成的`cert.pem`和`key.pem`拷贝到目录下面

```
$ tree
.
├── cert.pem
├── go.mod
├── go.sum
├── key.pem
└── main.go
```

```go
package main

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

const idleTimeout = 5 * time.Minute
const activeTimeout = 10 * time.Minute

func main() {
	var srv http.Server
	//http2.VerboseLogs = true
	srv.Addr = ":8972"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello http2"))
	})
	http2.ConfigureServer(&srv, &http2.Server{})
	log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
}
```

运行程序

```bash
$ go run main.go
```

然后用chrome浏览器可以正常浏览 `https://localhost:8972/` 

我们使用curl工具测试下http2和http1 client的情况:

```bash
$ curl --http2 https://localhost:8972/ -v
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8972 (#0)
...
> GET / HTTP/2
> Host: localhost:8972
> User-Agent: curl/7.60.0
> Accept: */*
> 
* Connection state changed (MAX_CONCURRENT_STREAMS == 250)!
< HTTP/2 200 
< content-type: text/plain; charset=utf-8
< content-length: 11
< date: Thu, 12 Dec 2019 11:13:26 GMT
< 
* Connection #0 to host localhost left intact
hello http2
```

我们看到使用http2协议返回正常。

我们在看下使用http1的情况：

```bash
$curl  --http1.1 https://localhost:8972/ -v
*   Trying 127.0.0.1...
...
* SSL connection using TLSv1.2 / ECDHE-RSA-AES128-GCM-SHA256
* ALPN, server accepted to use http/1.1
* Server certificate:
*  subject: O=mkcert development certificate; OU=wida@wida
*  start date: Jun  1 00:00:00 2019 GMT
*  expire date: Dec 12 10:47:26 2029 GMT
*  subjectAltName: host "localhost" matched cert's "localhost"
*  issuer: O=mkcert development CA; OU=wida@wida; CN=mkcert wida@wida
*  SSL certificate verify ok.
> GET / HTTP/1.1
> Host: localhost:8972
> User-Agent: curl/7.60.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Date: Thu, 12 Dec 2019 11:56:46 GMT
< Content-Length: 11
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact
hello http2
```

正常返回内容`hello http2`但是我们看到的是走的http1.1的协议。