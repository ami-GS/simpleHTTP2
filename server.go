package main

import (
	"fmt"
	hpack "github.com/ami-GS/GoHPACK"
	"net"
)

func StartServer() {
	serv, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := serv.Accept()
		if err != nil {
			panic(err)
		}
		client := NewSession(conn)
		client.Run()
	}
}
