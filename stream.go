package http2

type Stream struct {
	ID         uint32
	WindowSize uint16
	State      STATE
}

// 65535 should be defined in connection
func NewStream(streamID uint32) *Stream {
	return &Stream{streamID, 65535, IDLE}
}

func (self *Stream) ChangeState(state STATE) {
	self.State = state
}
