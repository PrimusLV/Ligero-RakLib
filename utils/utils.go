package utils

import (
	"bytes"
	"log"
	"utils/binpacker"
)

var buffer bytes.Buffer = new(bytes.Buffer)

//
// Reads 2 byte big-endian number
//
func ReadShort(bin string) int16 {
	defer buffer.Reset()
	buffer.WriteString(bin)
	up := binpacker.NewUnpacker(buffer)
	var val int16
	val, err := up.ShiftInt16()
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func WriteShort(val int16) string {
	defer buffer.Reset()
	p := binpacker.NewPacker()
	p.PushInt16(val)
	return string(buffer.Bytes())
}

func ReadLong(bin string) int64 {
	defer buffer.Reset()
	buffer.WriteString(bin)
	up := binpacker.NewUnpacker(buffer)
	var val int64
	val, err := up.ShiftInt64()
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func ReadInt(bin string) int {
	defer buffer.Reset()
	buffer.WriteString(bin)
	up := binpacker.NewUnpacker(buffer)
	var val int32
	val, err := up.ShiftInt32()
	if err != nil {
		log.Fatal(err)
	}
	return int(val)
}

func WriteInt(val int) string {
	defer buffer.Reset()
	p := binpacker.NewPacker(buffer)
	p.PushInt32(val)
	return string(buffer.Bytes())
}
