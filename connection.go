package http2

import (
	"fmt"
	hpack "github.com/ami-GS/GoHPACK"
	"net"
	"reflect"
)

type Connection struct {
	Conn                 net.Conn
	Streams              map[uint32]*Stream
	lastStreamID         uint32
	Table                *hpack.Table
	HeaderTableSize      uint32
	EnablePush           byte
	MaxConcurrentStreams uint32
	InitialWindowSize    uint32
	MaxFrameSize         uint32
	MaxHeaderListSize    uint32
	Preface              bool
	buf                  []byte
}

func (self *Connection) Parse(info *Http2Header) (frame Frame) {

	switch info.Type {
	case DATA_FRAME:
		frame = &Data{Header: info}
	case HEADERS_FRAME:
		frame = &Headers{Header: info, Table: self.Table} // not cool using table here
	case PRIORITY_FRAME:
		frame = &Priority{Header: info}
	case RST_STREAM_FRAME:
		frame = &Rst_stream{Header: info}
	case SETTINGS_FRAME:
		frame = &Settings{Header: info}
	case PING_FRAME:
		frame = &Ping{Header: info}
	case GOAWAY_FRAME:
		frame = &GoAway{Header: info}
	case WINDOW_UPDATE_FRAME:
		frame = &WindowUpdate{Header: info}
	case CONTINUATION_FRAME:
		frame = &Continuation{Header: info}
	default:
		panic("undefined frame type")
	}
	self.buf = make([]byte, info.Length)
	self.Conn.Read(self.buf)
	frame.Parse(self.buf)

	return frame
}

func (self *Connection) Send(frame Frame) {
	self.Streams[frame.GetStreamID()].Send(frame)
}

func (self *Connection) RunReceiver() {
	for {
		if self.Preface {
			self.buf = make([]byte, 9)
			self.Conn.Read(self.buf)
			if reflect.DeepEqual(self.buf, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}) {
				self.Preface = false
				return // connection closed from client
			}
			info := Http2Header{}
			info.Parse(self.buf)
			ID := info.GetStreamID()
			_, ok := self.Streams[ID]
			if !ok {
				// not cool
				self.AddStream(ID)
			}
			frame := self.Parse(&info)
			fmt.Printf("Receive: %s\n%s\n", self.Streams[ID].String(), frame.String())
			self.Streams[ID].EvaluateFrame(frame)
		} else {
			self.buf = make([]byte, 24)
			_, err := self.Conn.Read(self.buf)
			if err != nil {
				return
			} else {
				if reflect.DeepEqual(self.buf, CONNECTION_PREFACE) {
					self.Preface = true
					fmt.Printf("New connection from %v\n", self.Conn.RemoteAddr())
				}
			}
		}
	}
}

func (self *Connection) AddStream(streamID uint32) {
	self.Streams[streamID] = NewStream(self, streamID)
}

func (self *Connection) SetHeaderTableSize(value uint32) {
	self.HeaderTableSize = value
	self.Table.SetDynamicTableSize(value)
}

func NewConnection(conn net.Conn, streamID uint32) *Connection {
	table := hpack.NewTable()
	connection := Connection{conn, nil, 0, &table, 4096, 1, INFINITE, 65535, MAX_FRAME_SIZE_MIN, INFINITE, false, []byte{}}
	connection.Streams = map[uint32]*Stream{0: NewStream(&connection, 0)}
	connection.AddStream(streamID)
	return &connection
}
