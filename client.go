package http2

import (
	"fmt"
	"net"
)

func Connect(addr string) (client *Connection) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	client = NewConnection(conn, 1)
	client.Conn.Write(CONNECTION_PREFACE)
	fmt.Printf("Successflly connected to %s\n", addr)

	return client
}
