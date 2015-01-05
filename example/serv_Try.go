package main

import (
	http2 "github.com/ami-GS/simpleHTTP2"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 {
		http2.StartServer("127.0.0.1:" + args[0])
	} else {
		http2.StartServer("127.0.0.1:8080")
	}
}
