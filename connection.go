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

func (self *Connection) GetFrame(info *Http2Header) (frame Frame) {

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

	return frame
}

func (self *Connection) Send(frame Frame) {
	self.Streams[frame.GetStreamID()].Send(frame)
}

func (self *Connection) Recv(length uint32) (buffer []byte, err error) {
	buffer = make([]byte, length)
	_, err = self.Conn.Read(buffer)
	return
}

func (self *Connection) RunReceiver() {
	var buffer []byte
	var err error // not cool
	for {
		if self.Preface {
			buffer, err = self.Recv(9)
			if err != nil {
				self.Preface = false
				return
			}
			info := Http2Header{}
			info.Parse(buffer)
			ID := info.GetStreamID()
			_, ok := self.Streams[ID]
			if !ok {
				// not cool
				self.AddStream(ID)
			}
			frame := self.GetFrame(&info)
			buffer, err = self.Recv(info.Length)
			frame.Parse(buffer)
			fmt.Printf("%s: \t%s\n\t%s\n", RecvC.Apply("Receive"),
				self.Streams[ID].String(), frame.String())
			self.Streams[ID].EvaluateFrame(frame)
		} else {
			buffer, err = self.Recv(24)
			if reflect.DeepEqual(buffer, CONNECTION_PREFACE) {
				self.Preface = true
				fmt.Printf("New connection from %v\n", self.Conn.RemoteAddr())
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
