package http2

import (
	"fmt"
	//hpack "github.com/ami-GS/GoHPACK"
	"net"
)

func StartServer(addr string) {
	serv, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server is now running at %s\n", addr)

	for {
		conn, err := serv.Accept()
		if err != nil {
			panic(err)
		}
		client := NewSession(conn)
		client.RunReceiver()
	}
}
