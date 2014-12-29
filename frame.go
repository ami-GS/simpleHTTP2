package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	hpack "github.com/ami-GS/GoHPACK"
)

func http2Frame(length uint32, frame, flag byte, streamID uint32) []byte {
	header := make([]byte, 9)
	for i := 0; i < 3; i++ {
		header[i] = byte(length >> (byte(2-i) * 8))
	}
	header[3], header[4] = frame, flag
	for i := 0; i < 4; i++ {
		header[i+5] = byte(streamID >> (byte(3-i) * 8))
	}
	return header
}

func Data(data *string, flag, padLen byte) []byte {
	var frame []byte
	idx := 0
	if flag == FLAG_PADDED {
		frame = make([]byte, len(*data)+int(padLen+1))
		frame[idx] = padLen
		idx++
	} else {
		frame = make([]byte, uint32(len(*data)))
	}
	byteData := []byte(*data)
	for i, d := range byteData {
		frame[idx+i] = d
	}
	return frame
}

func Settings(streamID uint16, value uint32) []byte {
	frame := make([]byte, 6)
	for i := 0; i < 2; i++ {
		frame[i] = byte(streamID >> (byte(1-i) * 8))
	}
	for i := 0; i < 4; i++ {
		frame[i+2] = byte(value >> (byte(3-i) * 8))
	}
	return frame
}

func Headers(headers []hpack.Header, flag, padLen, weight byte, streamDependency uint32, table *hpack.Table) []byte {
	var frame []byte
	idx := 0
	wire, err := hex.DecodeString(hpack.Encode(headers, false, false, false, table, -1))
	fmt.Println(wire)
	if err != nil {
		panic(err)
	}
	if flag == FLAG_PADDED {
		frame = make([]byte, int(padLen+1)+len(wire))
		frame[idx] = padLen
		idx++
	} else if flag == FLAG_PRIORITY {
		frame = make([]byte, 5+len(wire))
		for i := 0; i < 4; i++ {
			frame[i] = byte(streamDependency >> (byte(3-i) * 8))
		}
		frame[4] = weight
		idx = 5
	} else if flag == FLAG_END_HEADERS || flag == FLAG_END_STREAM {
		frame = make([]byte, len(wire))
	} else {
		panic("undefined flag")
	}
	for i, w := range wire {
		frame[idx+i] = w
	}
	return frame
}

func main() {
	table := hpack.InitTable()
	headers := []hpack.Header{hpack.Header{":method", "GET"}, hpack.Header{":scheme", "http"},
		hpack.Header{":authority", "127.0.0.1"}, hpack.Header{":path", "/"}}
	hh := Headers(headers, FLAG_PRIORITY, 0, 255, 0xef00ff00, &table)
	fmt.Println(hh)
	//a := "string"
	aa := uint32((1 << 31) - 1)
	bi := make([]byte, 10)
	//	b := []byte(a)
	binary.LittleEndian.PutUint32(bi, aa)
	//	fmt.Println(bi)
	//	fmt.Println(a)
	//	fmt.Println(hex.DecodeString("\x00"))
	//	fmt.Println(hex.Dump(b))
}
