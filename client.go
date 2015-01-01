package http2

import (
	"net"
)

func Connect(addr string) (client Session) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	client = NewSession(conn)
	return client
}
