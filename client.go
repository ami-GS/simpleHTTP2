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
	client = NewSession(conn)
	client.Send(CONNECTION_PREFACE)
	fmt.Printf("Successflly connected to %s\n", addr)

	return client
}
