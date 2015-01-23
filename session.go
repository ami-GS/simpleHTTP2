package http2

import (
	"fmt"
	hpack "github.com/ami-GS/GoHPACK"
	"net"
	"reflect"
)

type Session struct {
	Conn  net.Conn
	Table hpack.Table
}

func (self *Session) Parse(buf []byte) {
	info := Http2Header{}
	info.Parse(buf[:9])

	var frame Frame
	if info.Type == TYPE_DATA {
		frame = &Data{Header: &info}
		frame.Parse(buf[9:])
	} else if info.Type == TYPE_HEADERS {
		var idx, padLen byte = 0, 0
		if info.Flag&FLAG_PADDED == FLAG_PADDED {
			padLen = buf[9]
			idx += 1 + padLen
		}
		if info.Flag&FLAG_PRIORITY == FLAG_PRIORITY {
			idx += 5
		}
		header := hpack.Decode(buf[idx:byte(len(buf[9:]))-padLen], &self.Table)
		frame = &Headers{Header: &info, Headers: header}

		frame.Parse(buf[9:])
		if info.Flag == FLAG_END_HEADERS {
			self.Send(NewData("Hello! DATA frame", 1, FLAG_PADDED, 5))
		}
	} else if info.Type == TYPE_PRIORITY {
	} else if info.Type == TYPE_RST_STREAM {
	} else if info.Type == TYPE_SETTINGS {
		frame = &Settings{Header: &info}
		frame.Parse(buf[9:])
		if info.Flag == FLAG_NO {
			self.Send(NewSettings(SETTINGS_NO, 0, FLAG_ACK))
		} else if info.Flag == FLAG_ACK {
		}
	} else if info.Type == TYPE_PING {
	} else if info.Type == TYPE_GOAWAY {
		frame = &GoAway{Header: &info}
		frame.Parse(buf[9:])
	} else if info.Type == TYPE_WINDOW_UPDATE {
	} else if info.Type == TYPE_CONTINUATION {
	} else {
		panic("undefined frame type")
	}
	fmt.Printf("Receive: %s\n", frame.String())
}

func (self *Session) Send(frame Frame) {
	fmt.Printf("Send: %s\n", frame.String())
	self.Conn.Write(frame.GetWire())
}

func (self *Session) RunReceiver() {
	var buf []byte
	for {
		buf = make([]byte, 1024)
		_, err := self.Conn.Read(buf)
		if err != nil {
			return //EOF?
		} else {
			if reflect.DeepEqual(buf[:24], CONNECTION_PREFACE) {
				fmt.Printf("New connection from %v\n", self.Conn.RemoteAddr())
				continue
			}
			self.Parse(buf)
		}
	}
}

func NewSession(conn net.Conn) (client Session) {
	client.Conn = conn
	client.Table = hpack.InitTable()
	return
}
