package main

import (
	hpack "github.com/ami-GS/GoHPACK"
	http2 "github.com/ami-GS/simpleHTTP2"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]
	var client http2.Session
	if len(args) == 2 {
		client = http2.Connect(args[0] + ":" + args[1])
	} else {
		client = http2.Connect("127.0.0.1:8080")
	}
	go client.RunReceiver()
	client.Send(http2.NewSettings(http2.NO_SETTING, 0, http2.NO))
	time.Sleep(time.Second)
	headers := []hpack.Header{hpack.Header{":method", "GET"}, hpack.Header{":scheme", "http"},
		hpack.Header{":authority", "127.0.0.1"}, hpack.Header{":path", "/"}}
	client.Send(http2.NewHeaders(headers, &client.Table, 1, http2.END_HEADERS, 0, 0, false, 0))
	time.Sleep(time.Second)
	client.Send(http2.NewGoAway(1, http2.NO_ERROR, "DEBUG string!!"))
}
