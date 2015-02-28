package http2

type Stream struct {
	state STATE
}

func NewStream() *Stream {
	return &Stream{IDLE}
}
