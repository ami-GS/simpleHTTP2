package main

import (
	"fmt"
)

func parseHttp2Header(data *[]byte) (uint32, byte, byte, uint32) {
	length := uint32((*data)[0])<<16 | uint32((*data)[1])<<8 | uint32((*data)[2])
	streamID := uint32((*data)[5])<<24 | uint32((*data)[6])<<16 | uint32((*data)[7])<<8 | uint32((*data)[8])
	return length, (*data)[3], (*data)[4], streamID
}

func main() {
	fmt.Println(parseHttp2Header(&[]byte{0xff, 0xff, 0xff, 0x01, 0x02, 0xff, 0xff, 0xff, 0xff}))
}
