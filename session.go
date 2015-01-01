package main

import (
	"fmt"
	hpack "github.com/ami-GS/GoHPACK"
	"net"
)

type Session struct {
	Conn  net.Conn
	Table hpack.Table
}

func (self *Session) Parse(buf []byte) {
	info := Http2Header{}
	info.Parse(buf[:9])

	if info.Type == TYPE_DATA {
		data := Data{}
		data.Parse(buf[9:], info.Flag, info.Length)
	} else if info.Type == TYPE_HEADERS {
		headers := Headers{}
		headers.Parse(buf[9:], info.Flag, &self.Table)
	} else if info.Type == TYPE_PRIORITY {
		fmt.Println("priority")
	} else if info.Type == TYPE_RST_STREAM {
		fmt.Println("rst stream")
	} else if info.Type == TYPE_SETTINGS {
		settings := Settings{}
		settings.Parse(buf[9:], info.Flag)
	} else if info.Type == TYPE_PING {
		fmt.Println("ping")
	} else if info.Type == TYPE_GOAWAY {
		goaway := GoAway{}
		goaway.Parse(buf[9:])
	} else if info.Type == TYPE_WINDOW_UPDATE {
		fmt.Println("window update")
	} else if info.Type == TYPE_CONTINUATION {
		fmt.Println("continuation")
	} else {
		panic("undefined frame type")
	}
}

func (self *Session) Send(data []byte) {
	self.Conn.Write(data)
}

func (self *Session) RunReceiver() {
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

func NewSession(conn net.Conn) (client Session) {
	client.Conn = conn
	client.Table = hpack.InitTable()
	return
}
