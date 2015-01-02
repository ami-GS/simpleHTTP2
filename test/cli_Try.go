package main

import (
	hpack "github.com/ami-GS/GoHPACK"
	"http2"
	"time"
)

func main() {
	client := http2.Connect("127.0.0.1:8080")
	go client.RunReceiver()
	client.Send(http2.NewSettings(http2.SETTINGS_NO, 0, http2.FLAG_NO))
	time.Sleep(time.Second)
	headers := []hpack.Header{hpack.Header{":method", "GET"}, hpack.Header{":scheme", "http"},
		hpack.Header{":authority", "127.0.0.1"}, hpack.Header{":path", "/"}}
	client.Send(http2.NewHeaders(headers, &client.Table, 1, http2.FLAG_END_HEADERS, 0, 0, false, 0))
	time.Sleep(time.Second)
	client.Send(http2.NewGoAway(1, http2.NO_ERROR, "DEBUG string!!"))
}
