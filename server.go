package main

import (
	"fmt"
	hpack "github.com/ami-GS/GoHPACK"
	"net"
)

type Client struct {
	Conn  net.Conn
	Table hpack.Table
}

func (self *Client) Parse(buf []byte) {
	info := Http2Header{}
	info.Parse(buf[:9])

	if info.Type == TYPE_DATA {
		data := Data{}
		data.Parse(buf[9:], info.Flag, info.Length)
	} else if info.Type == TYPE_HEADERS {
		headers := Headers{}
		headers.Parse(buf[:9], info.Flag, &self.Table)
	} else if info.Type == TYPE_PRIORITY {
		fmt.Println("priority")
	} else if info.Type == TYPE_RST_STREAM {
		fmt.Println("rst stream")
	} else if info.Type == TYPE_SETTINGS {
		settings := Settings{}
		settings.Parse(buf[:9], info.Flag)
	} else if info.Type == TYPE_PING {
		fmt.Println("ping")
	} else if info.Type == TYPE_GOAWAY {
		goaway := GoAway{}
		goaway.Parse(buf[:9])
	} else if info.Type == TYPE_WINDOW_UPDATE {
		fmt.Println("window update")
	} else if info.Type == TYPE_CONTINUATION {
		fmt.Println("continuation")
	} else {
		panic("undefined frame type")
	}
}

func (self *Client) Run() {
	var buf []byte
	for {
		buf = make([]byte, 1024)
		_, err := self.Conn.Read(buf)
		if err != nil {
			panic(err)
		}
		self.Parse(buf)
	}
}

func NewClient(conn net.Conn) (client Client) {
	client.Table = hpack.InitTable()
	return
}

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
		client := NewClient(conn)
		client.Run()
	}
}
