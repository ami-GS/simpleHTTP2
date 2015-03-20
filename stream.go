package http2

import (
	"fmt"
)

type Stream struct {
	Conn       *Connection
	ID         uint32
	WindowSize uint16
	State      STATE
}

// 65535 should be defined in connection
func NewStream(connection *Connection, streamID uint32) *Stream {
	return &Stream{connection, streamID, 65535, IDLE}
}

func (self *Stream) ChangeState(state STATE) {
	self.State = state
}

func (self *Stream) GetState() STATE {
	return self.State
}

func (self *Stream) DecreaseWindow(size uint16) {
	self.WindowSize -= size
}

func (self *Stream) Send(frame Frame) {
	// do something to self
	fmt.Printf("Send: \n%s\n", frame.String())
	(*self.Conn).Conn.Write(frame.GetWire())
}

func (self *Stream) EvaluateFrame(frame Frame) {
	frame.Evaluate(*self)
	// do something to self, also error handling
}

func (self *Stream) String() string {
	return fmt.Sprintf("Stream: ID=%d, Status=%s, WindowSize=%d", self.ID, self.State.String(), self.WindowSize)
}
