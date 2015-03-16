package http2

import (
	"fmt"
	hpack "github.com/ami-GS/GoHPACK"
	"net"
	"reflect"
)

type Connection struct {
	Conn         net.Conn
	Streams      map[uint32]*Stream
	lastStreamID uint32
	Table        *hpack.Table
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
	case HEADERS_FRAME:
		frame = &Headers{Header: &info}
	case PRIORITY_FRAME:
		frame = &Priority{Header: &info}
	case RST_STREAM_FRAME:
		frame = &Rst_stream{Header: &info}
	case SETTINGS_FRAME:
		frame = &Settings{Header: &info}
	case PING_FRAME:
		frame = &Ping{Header: &info}
	case GOAWAY_FRAME:
		frame = &GoAway{Header: &info}
	case WINDOW_UPDATE_FRAME:
		frame = &WindowUpdate{Header: &info}
	case CONTINUATION_FRAME:
		frame = &Continuation{Header: &info}
	default:
		panic("undefined frame type")
	}
	frame.Parse(buf[9:])

	fmt.Printf("Receive: \n%s\n", frame.String())
	self.Streams[ID].EvaluateFrame(frame)
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
	connection := Connection{conn, nil, 0, &table}
	connection.Streams = map[uint32]*Stream{0: NewStream(&connection, 0)}
	connection.AddStream(streamID)
	return &connection
}
