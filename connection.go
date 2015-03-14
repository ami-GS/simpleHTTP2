package http2

import (
	"fmt"
	hpack "github.com/ami-GS/GoHPACK"
	"net"
	"reflect"
)

type Connection struct {
	Conn    net.Conn
	Streams map[uint32]*Stream
	Table   *hpack.Table
}

func (self *Connection) Parse(buf []byte) {
	info := Http2Header{}
	info.Parse(buf[:9])

	ID := info.GetStreamID()
	_, ok := self.Streams[ID]
	if !ok {
		// not cool
		self.AddStream(ID)
	}

	var frame Frame
	switch info.Type {
	case DATA_FRAME:
		frame = &Data{Header: &info}
		frame.Parse(buf[9:])
	case HEADERS_FRAME:
		frame = &Headers{Header: &info}
		frame.Parse(buf[9:])
	case PRIORITY_FRAME:
		frame = &Priority{Header: &info}
		frame.Parse(buf[9:])
	case RST_STREAM_FRAME:
	case SETTINGS_FRAME:
		frame = &Settings{Header: &info}
		frame.Parse(buf[9:])
	case PING_FRAME:
		frame = &Ping{Header: &info}
		frame.Parse(buf[9:])
	case GOAWAY_FRAME:
		frame = &GoAway{Header: &info}
		frame.Parse(buf[9:])
	case WINDOW_UPDATE_FRAME:
	case CONTINUATION_FRAME:
	default:
		panic("undefined frame type")
	}

	self.Streams[ID].EvaluateFrame(frame)
	fmt.Printf("Receive: \n%s\n", frame.String())
}

func (self *Connection) Send(frame Frame) {
	self.Streams[frame.GetStreamID()].Send(frame)
}

func (self *Connection) RunReceiver() {
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

func (self *Connection) AddStream(streamID uint32) {
	self.Streams[streamID] = NewStream(self, streamID)
}

func NewConnection(conn net.Conn, streamID uint32) *Connection {
	table := hpack.InitTable()
	connection := Connection{conn, nil, &table}
	connection.Streams = map[uint32]*Stream{0: NewStream(&connection, 0)}
	connection.AddStream(streamID)
	return &connection
}
