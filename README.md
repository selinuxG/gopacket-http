# gopacket-http

<p>
<a href="https://github.com/shixiaofeia/gopacket-http">
    <img src="https://badgen.net/badge/Github/gopacket-http?icon=github" alt="">
</a>
<a href="https://github.com/shixiaofeia/gopacket-http/LICENSE">
    <img alt="GitHub" src="https://img.shields.io/github/license/shixiaofeia/gopacket-http?style=flat-square">
</a>
<img src="https://img.shields.io/github/go-mod/go-version/shixiaofeia/gopacket-http.svg?style=flat-square" alt="">
<img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/shixiaofeia/gopacket-http?style=flat-square">
<img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/shixiaofeia/gopacket-http?style=social">
</p>

监听网卡流量, 过滤并组装HTTP请求和响应, 供旁路分析, 抓包等用途

参考项目 [netgraph](https://github.com/ga0/netgraph)

## 使用

1. 安装libpcap-dev 和 gcc

```sh
# Ubuntu
sudo apt install -y libpcap-dev gcc

# CentOS
sudo yum install -y libpcap-devel gcc

# MacOS(Homebrew)
brew install libpcap

```

2. 安装gopacket-http

```sh
go get -u github.com/shixiaofeia/gopacket-http
```

3. 在代码中导入

```go
import "github.com/shixiaofeia/gopacket-http/packet"
```

## 快速开始

```go
package main

import (
	"context"
	"github.com/shixiaofeia/gopacket-http/packet"
	"log"
)

var eventCh = make(chan interface{}, 1024)

func main() {
	go handle()
	if err := packet.NewPacketHandle(context.Background(), "en0", eventCh).Listen(); err != nil {
		log.Println(err.Error())
	}
}

func handle() {
	for i := range eventCh {
		data := i.(packet.Event)
		log.Printf("request uri: %s, response status: %v", data.Req.RequestURI, data.Resp.Status)
	}
}

```


## 可配置的参数

```go
package main

import (
	"context"
	"github.com/shixiaofeia/gopacket-http/packet"
	"log"
	"net/http"
	"time"
)

var (
	eventCh     = make(chan interface{}, 1024)
	ctx, cancel = context.WithCancel(context.Background())
)

func main() {
	go shutdown()
	srv := packet.NewPacketHandle(ctx, "en0", eventCh)
	srv.SetBpf("tcp port 80")     // 设置BPF过滤规则
	srv.SetEventHandle(5, handle) // 设置多协程事件处理,
	srv.SetPromisc(true)          // 设置混杂模式开启状态,
	srv.SetFlushTime(time.Minute) // 设置清理缓存时间
	if err := srv.Listen(); err != nil {
		log.Println(err.Error())
	}
}

func handle(req *http.Request, resp *http.Response) {
	log.Printf("request uri: %s, response status: %v", req.RequestURI, resp.Status)
}

func shutdown() {
	time.Sleep(time.Second * 10)
	cancel()
}

```


