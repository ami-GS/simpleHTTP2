package http2

import (
	"fmt"
	hpack "github.com/ami-GS/GoHPACK"
)

type Frame interface {
	Pack()
	Parse(data []byte)
	String() string
	GetWire() []byte
	GetStreamID() uint32
	Evaluate(Stream)
}

type Http2Header struct {
	Length   uint32
	Type     TYPE
	Flag     FLAG
	StreamID uint32
	HeadWire []byte
}

func NewHttp2Header(length uint32, fType TYPE, flag FLAG, streamID uint32) *Http2Header {
	h := Http2Header{length, fType, flag, streamID, []byte{}}
	h.Pack()
	return &h
}

func (self *Http2Header) Pack() {
	self.HeadWire = make([]byte, 9)
	for i := 0; i < 3; i++ {
		self.HeadWire[i] = byte(self.Length >> byte((2-i)*8))
	}
	self.HeadWire[3], self.HeadWire[4] = byte(self.Type), byte(self.Flag)
	for i := 0; i < 4; i++ {
		self.HeadWire[i+5] = byte(self.StreamID >> byte((3-i)*8))
	}
}

func (self *Http2Header) Parse(data []byte) {
	self.Length = uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
	self.Type = TYPE(data[3])
	self.Flag = FLAG(data[4])
	self.StreamID = uint32(data[5])<<24 | uint32(data[6])<<16 | uint32(data[7])<<8 | uint32(data[8])
}

func (self *Http2Header) String() string {
	str := fmt.Sprintf("%s frame: Length=%d, Flag=%s, StreamID=%d",
		self.Type.String(), self.Length, self.Flag.String(), self.StreamID)
	return str
}

func (self *Http2Header) GetWire() []byte {
	return self.HeadWire
}

func (self *Http2Header) GetStreamID() uint32 {
	return self.StreamID
}

type Data struct {
	Header *Http2Header
	Data   string
	PadLen byte
	Wire   []byte
}

func NewData(data string, streamID uint32, flag FLAG, padLen byte) *Data {
	var length uint32 = uint32(len(data))
	if flag&PADDED == PADDED {
		length += uint32(padLen + 1)
	}

	header := NewHttp2Header(length, DATA_FRAME, flag, streamID)
	frame := Data{header, data, padLen, []byte{}}
	frame.Pack()

	return &frame
}

func (self *Data) Pack() {
	idx := 0
	if self.Header.Flag == PADDED {
		self.Wire = make([]byte, len(self.Data)+int(self.PadLen+1))
		self.Wire[idx] = self.PadLen
		idx++
	} else {
		self.Wire = make([]byte, uint32(len(self.Data)))
	}
	for i, d := range self.Data {
		self.Wire[idx+i] = byte(d)
	}
}

func (self *Data) Parse(data []byte) {
	if self.Header.Flag == PADDED {
		self.PadLen = data[0]
		self.Data = string(data[1 : self.Header.Length-uint32(self.PadLen)])
	} else {
		self.Data = string(data)
	}
}

func (self *Data) String() string {
	return fmt.Sprintf("%s\n{contents:%s}", self.Header.String(), self.Data)
}

func (self *Data) GetWire() []byte {
	return append(self.Header.GetWire(), self.Wire...)
}

func (self *Data) GetStreamID() uint32 {
	return self.Header.GetStreamID()
}

func (self *Data) Evaluate(stream Stream) {
	if self.GetStreamID() == 0 {
		//stream.Send(NewGoAway(stream.lastID, PROTOCOL_ERROR, ""))
	}

	state := stream.GetState()
	if state == CLOSED {
		stream.Send(NewRst_stream(PROTOCOL_ERROR, stream.ID))
	}
	if state != OPEN && state != HALF_CLOSED_LOCAL {
		stream.Send(NewRst_stream(STREAM_CLOSED, stream.ID))
	}

	/*
		if padLen > (len) {
			stream.Send(NewGoAway(straem.lastID, PROTOCOL_ERROR, ""))
		}
	*/

	if self.Header.Flag&END_STREAM == END_STREAM {
		if state == OPEN {
			stream.ChangeState(HALF_CLOSED_REMOTE)
		} else if state == HALF_CLOSED_LOCAL {
			stream.ChangeState(CLOSED)
		}
	}

	//TODO: decrease window based on the data
}

type Settings struct {
	Header    *Http2Header
	SettingID SETTING
	Value     uint32
	Wire      []byte
}

func NewSettings(settingID SETTING, value uint32, flag FLAG) *Settings {
	header := NewHttp2Header(uint32(6), SETTINGS_FRAME, flag, 0)
	frame := Settings{header, settingID, value, []byte{}}
	frame.Pack()

	return &frame
}

func (self *Settings) Pack() {
	self.Wire = make([]byte, 6)
	for i := 0; i < 2; i++ {
		self.Wire[i] = byte(self.SettingID >> byte((1-i)*8))
	}
	for i := 0; i < 4; i++ {
		self.Wire[2+i] = byte(self.Value >> byte((3-i)*8))
	}
}

func (self *Settings) Parse(data []byte) {
	self.SettingID = SETTING(uint16(data[0])<<8 | uint16(data[1]))
	self.Value = uint32(data[2])<<24 | uint32(data[3])<<16 | uint32(data[4])<<8 | uint32(data[5])
	_ = self.Header.Flag //temporally
}

func (self *Settings) String() string {
	return fmt.Sprintf("%s\nsetting=%s(%d)",
		self.Header.String(), self.SettingID.String(), self.Value)
}

func (self *Settings) GetWire() []byte {
	return append(self.Header.GetWire(), self.Wire...)
}

func (self *Settings) GetStreamID() uint32 {
	return self.Header.GetStreamID()
}

func (self *Settings) Evaluate(stream Stream) {
	if stream.ID != 0 {
		//stream.Send(NewGoAway(stream.lastID, PROTOCOL_ERROR, ""))
	}
	if self.Header.Length%6 != 0 {
		//stream.Send(NewGoAway(stream.lastID, FRAME_SIZE_ERROR, ""))
	}
	if self.Header.Flag == ACK {
		if self.Header.Length != 0 {
			//stream.Send(NewGoAway(stream.lastID, FRAME_SIZE_ERROR, ""))
		}
	} else if self.Header.Length > 0 {
		if self.SettingID == HEADER_TABLE_SIZE {
			// setTableSize
		} else if self.SettingID == ENABLE_PUSH {
			if self.Value == 1 || self.Value == 0 {
				// setPush
			} else {
				//stream.Send(NewGoAway(stream.lastID, PROTOCOL_ERROR, ""))
			}
		} else if self.SettingID == MAX_CONCURRENT_STREAMS {
			if self.Value <= 100 {
				fmt.Println("Warnnig: max_concurrent_stream below 100 is not recomended")
			}
			// setMaxConcurrentStream
		} else if self.SettingID == INITIAL_WINDOW_SIZE {
			/*
				if self.Value > MAX_WINDOW_SIZE {
					stream.Send(NewGoAway(stream.lastID, FLOW_CONTROL_ERROR, ""))
				} else {
					// setInitialWindowSize
				}
			*/
		} else if self.SettingID == MAX_FRAME_SIZE {
			/*
				if INITIAL_MAX_FRAME_SIZE <= self.Value && self.Value <= LIMIT_MAX_FRAME_SIZE {
					// setMaxFrameSize
				} else {
					//stream.Send(NewGoAway(stream.lastID, PROTOCOL_ERROR, ""))
				}
			*/
		} else if self.SettingID == MAX_HEADER_LIST_SIZE {
			//setMaxHeaderListize
		} else {
			// ignore
		}
		stream.Send(NewSettings(NO_SETTING, 0, ACK))
	}
}

type Headers struct {
	Header           *Http2Header
	Headers          []hpack.Header
	block            []byte
	PadLen, Weight   byte
	E                bool
	StreamDependency uint32
	Wire             []byte
}

func NewHeaders(headers []hpack.Header, table *hpack.Table, streamID uint32, flag FLAG, padLen, weight byte, e bool, streamDependency uint32) *Headers {
	encHeaders := hpack.Encode(headers, false, false, false, table, -1)
	var length uint32 = uint32(len(encHeaders))
	if flag&PADDED == PADDED {
		length += uint32(padLen + 1)
	}
	if flag&PRIORITY == PRIORITY {
		length += 5
	}

	header := NewHttp2Header(length, HEADERS_FRAME, flag, streamID)

	frame := Headers{header, headers, encHeaders, padLen, weight, e, streamDependency, []byte{}}
	frame.Pack()

	return &frame
}

func (self *Headers) Pack() {
	idx := 0
	if self.Header.Flag&PADDED == PADDED {
		self.Wire = make([]byte, int(self.PadLen+1)+len(self.block))
		self.Wire[idx] = self.PadLen
		idx++
	}
	if self.Header.Flag&PRIORITY == PRIORITY {
		self.Wire = make([]byte, 5+len(self.block))
		for i := 0; i < 4; i++ {
			self.Wire[i] = byte(self.StreamDependency >> byte((3-i)*8))
		}
		if self.E {
			self.Wire[0] |= 0x80
		}
		self.Wire[4] = self.Weight
		idx = 5
	}
	if self.Header.Flag&END_HEADERS == END_HEADERS || self.Header.Flag&END_STREAM == END_STREAM {
		self.Wire = make([]byte, len(self.block))
	}
	/*else {
		panic("undefined flag")
	}*/
	for i, h := range self.block {
		self.Wire[idx+i] = h
	}
}

func (self *Headers) Parse(data []byte) {
	idx := 0
	if self.Header.Flag&PADDED == PADDED {
		self.PadLen = data[idx]
		idx++
	}
	if self.Header.Flag&PRIORITY == PRIORITY {
		if data[idx]&0x80 > 0 {
			self.E = true
		}
		self.StreamDependency = uint32(data[idx]&0xef)<<24 | uint32(data[idx+1])<<16 | uint32(data[idx+2])<<8 | uint32(data[idx+3])
		self.Weight = data[idx+4]
		idx += 5
	}
	/*else {
		panic("undefined flag")
	}*/
}

func (self *Headers) String() string {
	return fmt.Sprintf("%s\n{Headers:%v}", self.Header.String(), self.Headers)
}

func (self *Headers) GetWire() []byte {
	return append(self.Header.GetWire(), self.Wire...)
}

func (self *Headers) GetStreamID() uint32 {
	return self.Header.GetStreamID()
}

func (self *Headers) Evaluate(stream Stream) {
	if stream.GetState() == RESERVED_LOCAL {
		stream.ChangeState(HALF_CLOSED_LOCAL)
	} else {
		stream.ChangeState(OPEN)
	}

	if stream.ID == 0 {
		//stream.Send(NewGoAway(stream.lastID, PROTOCOL_ERROR, ""))
	}

	if self.Header.Flag&END_HEADERS == END_HEADERS {
		self.Headers = hpack.Decode(self.block, (*stream.Conn).Table)
		// The stream.ID is suspicious
		stream.Send(NewData("data:hoge", stream.ID, END_STREAM, 0))
		//stream.ChangeState(?)
	}
	if self.Header.Flag&END_STREAM == END_STREAM {
		stream.ChangeState(HALF_CLOSED_REMOTE)
	}
}

type Priority struct {
	Header           *Http2Header
	E                bool
	StreamDependency uint32
	Weight           byte
	Wire             []byte
}

func NewPriority(streamID uint32, e bool, streamDependency uint32, weight byte) *Priority {
	header := NewHttp2Header(5, PRIORITY_FRAME, NO, streamID)
	frame := Priority{header, e, streamDependency, weight, []byte{}}
	frame.Pack()
	return &frame
}

func (self *Priority) Pack() {
	self.Wire = make([]byte, 5)
	for i := 0; i < 4; i++ {
		self.Wire[i] = byte(self.StreamDependency >> (byte(3-i) * 8))
	}
	if self.E {
		self.Wire[0] |= 0x80
	}
	self.Wire[4] = self.Weight
}

func (self *Priority) Parse(data []byte) {
	if data[0]&0x80 > 0 {
		self.E = true
	}
	self.StreamDependency = uint32(data[0]&0xef)<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
	self.Weight = data[4]
}

func (self *Priority) String() string {
	return fmt.Sprintf("%s\n{Priority: E=%t, StreamDependency=%b, Weight=%d}", self.Header.String(), self.E, self.StreamDependency, self.Weight)
}

func (self *Priority) GetWire() []byte {
	return append(self.Header.GetWire(), self.Wire...)
}

func (self *Priority) GetStreamID() uint32 {
	return self.Header.GetStreamID()
}

func (self *Priority) Evaluate(stream Stream) {}

type Ping struct {
	Header   *Http2Header
	PingData string
	Wire     []byte
}

func NewPing(pingData string, flag FLAG) *Ping {
	header := NewHttp2Header(8, PING_FRAME, flag, 0)
	frame := Ping{header, pingData, []byte{}}
	frame.Pack()
	return &frame
}

func (self *Ping) Pack() {
	self.Wire = make([]byte, 8)
	self.Wire = []byte(self.PingData)
}

func (self *Ping) Parse(data []byte) {
	self.PingData = string(data)
}

func (self *Ping) String() string {
	return fmt.Sprintf("%s\n{ping:%s}", self.Header.String(), self.PingData)
}

func (self *Ping) GetWire() []byte {
	return append(self.Header.GetWire(), self.Wire...)
}

func (self *Ping) GetStreamID() uint32 {
	return self.Header.GetStreamID()
}

func (self *Ping) Evaluate(stream Stream) {}

type Rst_stream struct {
	Header    *Http2Header
	ErrorCode ERROR
	Wire      []byte
}

func NewRst_stream(errorCode ERROR, streamID uint32) *Rst_stream {
	header := NewHttp2Header(4, RST_STREAM_FRAME, NO, streamID)
	frame := Rst_stream{header, errorCode, []byte{}}
	frame.Pack()
	return &frame
}

func (self *Rst_stream) Pack() {
	self.Wire = make([]byte, 4)
	for i := 0; i < 4; i++ {
		self.Wire[i] = byte(self.ErrorCode >> byte((3-i)*8))
	}
}

func (self *Rst_stream) Parse(data []byte) {
	self.ErrorCode = ERROR(uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3]))
}

func (self *Rst_stream) String() string {
	return fmt.Sprintf("%s\n{Error:%s}", self.Header.String(), self.ErrorCode.String())
}

func (self *Rst_stream) GetWire() []byte {
	return append(self.Header.GetWire(), self.Wire...)
}

func (self *Rst_stream) GetStreamID() uint32 {
	return self.Header.GetStreamID()
}

func (self *Rst_stream) Evaluate(stream Stream) {}

type GoAway struct {
	Header       *Http2Header
	LastStreamID uint32
	ErrorCode    ERROR
	Debug        string
	Wire         []byte
}

func NewGoAway(lastStreamID uint32, errorCode ERROR, debug string) *GoAway {
	var length uint32 = uint32(len(debug) + 8)
	header := NewHttp2Header(length, GOAWAY_FRAME, NO, 0)
	frame := GoAway{header, lastStreamID, errorCode, debug, []byte{}}
	frame.Pack()
	return &frame
}

func (self *GoAway) Pack() {
	self.Wire = make([]byte, 8+len(self.Debug))
	for i := 0; i < 4; i++ {
		self.Wire[i] = byte(self.LastStreamID >> byte((3-i)*8))
		self.Wire[i+4] = byte(self.ErrorCode >> byte((3-i)*8))
	}
	for i, d := range self.Debug {
		self.Wire[i+8] = byte(d)
	}
}

func (self *GoAway) Parse(data []byte) {
	self.LastStreamID = uint32(data[0]&0xef)<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
	self.ErrorCode = ERROR(uint32(data[4])<<24 | uint32(data[5])<<16 | uint32(data[6])<<8 | uint32(data[7]))
	if len(data) >= 9 {
		self.Debug = string(data[8:])
	}
}

func (self *GoAway) String() string {
	return fmt.Sprintf("%s\n{debug:%s}", self.Header.String(), self.Debug)
}

func (self *GoAway) GetWire() []byte {
	return append(self.Header.GetWire(), self.Wire...)
}

func (self *GoAway) GetStreamID() uint32 {
	return self.Header.GetStreamID()
}

func (self *GoAway) Evaluate(stream Stream) {}

type WindowUpdate struct {
	Header              *Http2Header
	WindowSizeIncrement uint32
	Wire                []byte
}

func NewWindowUpdate(streamID, windowSizeIncrement uint32) *WindowUpdate {
	header := NewHttp2Header(4, WINDOW_UPDATE_FRAME, NO, streamID)
	frame := WindowUpdate{header, windowSizeIncrement, []byte{}}
	frame.Pack()
	return &frame
}

func (self *WindowUpdate) Pack() {
	self.Wire = make([]byte, 4)
	self.Wire[0] = byte((self.WindowSizeIncrement >> 24) & 0xef) // top bit is reserved
	for i := 1; i < 4; i++ {
		self.Wire[i] = byte(self.WindowSizeIncrement >> (byte(3-i) * 8))
	}
}

func (self *WindowUpdate) Parse(data []byte) {
	self.WindowSizeIncrement = uint32(data[0]&0xef)<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
}

func (self *WindowUpdate) String() string {
	return fmt.Sprintf("%s\n{WindowUpdate: WindowSizeIncrement=%d}", self.Header.String(), self.WindowSizeIncrement)
}

func (self *WindowUpdate) GetWire() []byte {
	return append(self.Header.GetWire(), self.Wire...)
}

func (self *WindowUpdate) GetStreamID() uint32 {
	return self.Header.GetStreamID()
}

func (self *WindowUpdate) Evaluate(stream Stream) {}

type Continuation struct {
	Header *Http2Header
	Block  []byte
	Wire   []byte
}

func NewContinuation(block []byte, streamID uint32, flag FLAG) *Continuation {
	header := NewHttp2Header(uint32(len(block)), CONTINUATION_FRAME, flag, streamID)
	frame := Continuation{header, block, []byte{}}
	frame.Pack()
	return &frame
}

func (self *Continuation) Pack() {
	self.Wire = self.Block
}

func (self *Continuation) Parse(data []byte) {
	self.Wire = data
}

func (self *Continuation) String() string {
	return fmt.Sprintf("%s\n{Continuation}", self.Header.String())
}

func (self *Continuation) GetWire() []byte {
	return append(self.Header.GetWire(), self.Wire...)
}

func (self *Continuation) GetStreamID() uint32 {
	return self.Header.GetStreamID()
}

func (self *Continuation) Evaluate(stream Stream) {}

/*
func main() {
	table := hpack.InitTable()
	headers := []hpack.Header{hpack.Header{":method", "GET"}, hpack.Header{":scheme", "http"},
		hpack.Header{":authority", "127.0.0.1"}, hpack.Header{":path", "/"}}

	http2Header := Http2Header{Length: 12, Type: TYPE_DATA, Flag: FLAG_PADDED, StreamID: 1}
	http2Header.Pack()
	fmt.Printf("http2Header %v\n", http2Header)
	data := Data{Data: "Hello!", PadLen: 5}
	data.Pack(http2Header.Flag)
	data2 := Data{}
	data2.Parse(data.Wire, http2Header.Flag, http2Header.Length)
	fmt.Printf("data %v\n", data)
	fmt.Printf("data2 %v\n", data2)
	settings := Settings{SettingID: 0xff00, Value: 0xff00ff00}
	settings2 := Settings{}
	settings.Pack()
	settings2.Parse(settings.Wire, http2Header.Flag)
	fmt.Printf("settings %v\n", settings)
	fmt.Printf("settings2 %v\n", settings2)
	h := Headers{Headers: headers, PadLen: 5, Weight: 0, E: false}
	h2 := Headers{}
	h.Pack(http2Header.Flag, &table)
	h2.Parse(h.Wire, http2Header.Flag, &table)
	fmt.Printf("headers %v\n", h)
	fmt.Printf("headers2 %v\n", h2)
	goaway := GoAway{LastStreamID: 0xef00ff00, ErrorCode: 0xff00ff00, Debug: "DEBUG MESSAGE!!"}
	goaway2 := GoAway{}
	goaway.Pack()
	goaway2.Parse(goaway.Wire)
	fmt.Printf("goaway %v\n", goaway)
	fmt.Printf("goaway2 %v\n", goaway2)
}
*/
