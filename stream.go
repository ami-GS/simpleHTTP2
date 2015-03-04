package http2

type Stream struct {
	connection *Connection
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

func (self *Stream) DecreaseWindow(size uint16) {
	self.WindowSize -= size
}

func (self *Stream) Send(frame Frame) {
	self.connection.Send(frame)
}
