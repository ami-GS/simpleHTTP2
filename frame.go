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

/*func Data(string data, flag bool, uint32 padLen) byte {

	return
}*/

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
