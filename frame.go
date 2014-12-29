package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func http2Frame(length uint32, frame, flag byte, streamID uint32) []byte {
	header := make([]byte, 9)
	for i := 0; i < 3; i++ {
		header[i] = byte(length>>byte(2-i)) & 0xff
	}
	header[3], header[4] = frame, flag
	for i := 0; i < 4; i++ {
		header[i+5] = byte(streamID>>byte(3-i)) & 0xff
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

func main() {
	a := "string"
	aa := uint32((1 << 31) - 1)
	bi := make([]byte, 10)
	b := []byte(a)
	binary.LittleEndian.PutUint32(bi, aa)
	fmt.Println(bi)
	fmt.Println(a)
	fmt.Println(hex.DecodeString("\x00"))
	fmt.Println(hex.Dump(b))
}
