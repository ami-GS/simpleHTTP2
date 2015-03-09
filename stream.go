package http2

import (
	"fmt"
	"net"
)

type Stream struct {
	Conn       *net.Conn
	ID         uint32
	WindowSize uint16
	State      STATE
}

// 65535 should be defined in connection
func NewStream(conn *net.Conn, streamID uint32) *Stream {
	return &Stream{conn, streamID, 65535, IDLE}
}

func (self *Stream) ChangeState(state STATE) {
	self.State = state
}

func (self *Stream) DecreaseWindow(size uint16) {
	self.WindowSize -= size
}

func (self *Stream) Send(frame Frame) {
	// do something to self
	fmt.Printf("Send: \n%s\n", frame.String())
	(*self.Conn).Write(frame.GetWire())
}

func (self *Stream) EvaluateFrame(frame Frame) {
	// do something to self, also error handling
}
