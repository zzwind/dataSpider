package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

func main() {
	fmt.Println(Float32ToByte(0.00))
}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}
