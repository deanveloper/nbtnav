package nbt

import (
	"encoding/binary"
	"io"
	"bufio"
)

func Parse(r io.Reader) (*Tag, error) {
	// TODO
	return nil, nil
}

func ReadTagType(buf *bufio.Reader) (TagType, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return TagEnd, err
	}
	return TagType(b), nil
}

func ParseInt8(buf *bufio.Reader) (int8, error) {
	var i int8
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func ParseInt16(buf *bufio.Reader) (int16, error) {
	var i int16
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func ParseInt32(buf *bufio.Reader) (int32, error) {
	var i int32
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func ParseInt64(buf *bufio.Reader) (int64, error) {
	var i int64
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func ParseFloat32(buf *bufio.Reader) (float32, error) {
	var i float32
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func ParseFloat64(buf *bufio.Reader) (float64, error) {
	var i float64
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func ParseInt8Array(buf *bufio.Reader) ([]int8, error) {
	length, err := ParseInt32(buf)
	if err != nil {
		return nil, err
	}
	arr := make([]int8, length)
	for i := 0; i < int(length); i++ {
		i8, err := ParseInt8(buf)
		if err != nil {
			return nil, err
		}
		arr[i] = i8
	}
	return arr, nil
}

func ParseString(buf *bufio.Reader) (string, error) {
	length, err := ParseInt16(buf)
	if err != nil {
		return "", err
	}

	arr := make([]byte, length)
	for i := 0; i < int(length); i++ {
		i8, err := ParseInt8(buf)
		if err != nil {
			return "", err
		}
		arr[i] = byte(i8)
	}

	return string(arr), nil
}

// TODO add list, compound, intarray
