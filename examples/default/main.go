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
