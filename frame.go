package http2

import (
	"fmt"
	hpack "github.com/ami-GS/GoHPACK"
)

type Http2Header struct {
	Wire       []byte
	Length     uint32
	Type, Flag byte
	StreamID   uint32
}

func (self *Http2Header) Pack() {
	self.Wire = make([]byte, 9)
	for i := 0; i < 3; i++ {
		self.Wire[i] = byte(self.Length >> byte((2-i)*8))
	}
	self.Wire[3], self.Wire[4] = self.Type, self.Flag
	for i := 0; i < 4; i++ {
		self.Wire[i+5] = byte(self.StreamID >> byte((3-i)*8))
	}
}

func (self *Http2Header) Parse(data []byte) {
	self.Length = uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
	self.Type = data[3]
	self.Flag = data[4]
	self.StreamID = uint32(data[5])<<24 | uint32(data[6])<<16 | uint32(data[7])<<8 | uint32(data[8])
}

type Data struct {
	Wire   []byte
	Data   string
	PadLen byte
}

func NewData(data string, streamID uint32, flag, padLen byte) []byte {
	frame := Data{Data: data, PadLen: padLen}
	frame.Pack(flag)
	header := Http2Header{Length: uint32(len(frame.Wire)), Type: TYPE_DATA, Flag: flag, StreamID: streamID}
	header.Pack()
	return append(header.Wire, frame.Wire...)
}

func (self *Data) Pack(flag byte) {
	idx := 0
	if flag == FLAG_PADDED {
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

func (self *Data) Parse(data []byte, flag byte, length uint32) {
	if flag == FLAG_PADDED {
		self.PadLen = data[0]
		self.Data = string(data[1 : length-uint32(self.PadLen)])
	} else {
		self.Data = string(data)
	}
}

type Settings struct {
	Wire      []byte
	SettingID uint16
	Value     uint32
}

func NewSettings(settingID uint16, value uint32, flag byte) []byte {
	frame := Settings{SettingID: settingID, Value: value}
	frame.Pack()
	header := Http2Header{Length: uint32(len(frame.Wire)), Type: TYPE_SETTINGS, Flag: flag, StreamID: 0}
	header.Pack()
	return append(header.Wire, frame.Wire...)
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

func (self *Settings) Parse(data []byte, flag byte) {
	self.SettingID = uint16(data[0])<<8 | uint16(data[1])
	self.Value = uint32(data[2])<<24 | uint32(data[3])<<16 | uint32(data[4])<<8 | uint32(data[5])
	_ = flag //temporally
}

type Headers struct {
	Wire             []byte
	Headers          []hpack.Header
	PadLen, Weight   byte
	E                bool
	StreamDependency uint32
}

func NewHeaders(headers []hpack.Header, table *hpack.Table, streamID uint32, flags, padLen, weight byte, e bool, streamDependency uint32) []byte {
	frame := Headers{Headers: headers, PadLen: padLen, Weight: weight, E: e, StreamDependency: streamDependency}
	frame.Pack(flags, table)
	header := Http2Header{Length: uint32(len(frame.Wire)), Type: TYPE_HEADERS, Flag: flags, StreamID: streamID}
	header.Pack()
	return append(header.Wire, frame.Wire...)
}

func (self *Headers) Pack(flags byte, table *hpack.Table) {
	idx := 0
	encHeaders := hpack.Encode(self.Headers, false, false, false, table, -1)
	if flags&FLAG_PADDED == FLAG_PADDED {
		self.Wire = make([]byte, int(self.PadLen+1)+len(encHeaders))
		self.Wire[idx] = self.PadLen
		idx++
	}
	if flags&FLAG_PRIORITY == FLAG_PRIORITY {
		self.Wire = make([]byte, 5+len(encHeaders))
		for i := 0; i < 4; i++ {
			self.Wire[i] = byte(self.StreamDependency >> byte((3-i)*8))
		}
		if self.E {
			self.Wire[0] |= 0x80
		}
		self.Wire[4] = self.Weight
		idx = 5
	}
	if flags&FLAG_END_HEADERS == FLAG_END_HEADERS || flags&FLAG_END_STREAM == FLAG_END_STREAM {
		self.Wire = make([]byte, len(encHeaders))
	}
	/*else {
		panic("undefined flag")
	}*/
	for i, h := range encHeaders {
		self.Wire[idx+i] = h
	}
}

func (self *Headers) Parse(data []byte, flags byte, table *hpack.Table) {
	idx := 0
	if flags&FLAG_PADDED == FLAG_PADDED {
		self.PadLen = data[idx]
		idx++
	}
	if flags&FLAG_PRIORITY == FLAG_PRIORITY {
		if data[0]&0x80 > 0 {
			self.E = true
		}
		self.StreamDependency = uint32(data[0]&0xef)<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
		self.Weight = data[4]
		idx += 5
	}
	if flags&FLAG_END_HEADERS == FLAG_END_HEADERS || flags&FLAG_END_STREAM == FLAG_END_STREAM {
		fmt.Println("change stream state")
	}
	/*else {
		panic("undefined flag")
	}*/
	self.Headers = hpack.Decode(data[idx:len(data)-int(self.PadLen)], table)
}

type GoAway struct {
	Wire         []byte
	LastStreamID uint32
	ErrorCode    uint32
	Debug        string
}

func NewGoAway(lastStreamID, errorCode uint32, debug string) []byte {
	frame := GoAway{LastStreamID: lastStreamID, ErrorCode: errorCode, Debug: debug}
	frame.Pack()
	header := Http2Header{Length: uint32(len(frame.Wire)), Type: TYPE_GOAWAY, Flag: FLAG_NO, StreamID: 0}
	header.Pack()
	return append(header.Wire, frame.Wire...)

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
	self.ErrorCode = uint32(data[4])<<24 | uint32(data[5])<<16 | uint32(data[6])<<8 | uint32(data[7])
	if len(data) >= 9 {
		self.Debug = string(data[8:])
	}
}

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
