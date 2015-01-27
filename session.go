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
	if info.Type == DATA_FRAME {
		frame = &Data{Header: &info}
		frame.Parse(buf[9:])
	} else if info.Type == HEADERS_FRAME {
		var idx, padLen byte = 0, 0
		if info.Flag&PADDED == PADDED {
			padLen = buf[9]
			idx += 1
		}
		if info.Flag&PRIORITY == PRIORITY {
			idx += 5
		}
		header := hpack.Decode(buf[9+idx:info.Length-uint32(padLen)], &self.Table)
		frame = &Headers{Header: &info, Headers: header}

		frame.Parse(buf[9:])
		if info.Flag == END_HEADERS {
			self.Send(NewData("Hello! DATA frame", 1, PADDED, 5))
		}
	} else if info.Type == PRIORITY_FRAME {
	} else if info.Type == RST_STREAM_FRAME {
	} else if info.Type == SETTINGS_FRAME {
		frame = &Settings{Header: &info}
		frame.Parse(buf[9:])
		if info.Flag == NO {
			self.Send(NewSettings(NO_SETTING, 0, ACK))
		} else if info.Flag == ACK {
		}
	} else if info.Type == PING_FRAME {
	} else if info.Type == GOAWAY_FRAME {
		frame = &GoAway{Header: &info}
		frame.Parse(buf[9:])
	} else if info.Type == WINDOW_UPDATE_FRAME {
	} else if info.Type == CONTINUATION_FRAME {
	} else {
		panic("undefined frame type")
	}
	fmt.Printf("Receive: \n%s\n", frame.String())
}

func (self *Session) Send(frame Frame) {
	fmt.Printf("Send: \n%s\n", frame.String())
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

func NewSession(conn net.Conn) (session Session) {
	session.Conn = conn
	session.Table = hpack.InitTable()
	return
}
