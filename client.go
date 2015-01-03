package http2

import (
	"fmt"
	"net"
)

func Connect(addr string) (client Session) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successflly connected to %s\n", addr)

	client = NewSession(conn)
	return client
}
