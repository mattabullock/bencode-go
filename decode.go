package main

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
)

func ParseBencode(enc []byte) interface{} {
	byteReader := bytes.NewReader(enc)
	rdr := bufio.NewReader(byteReader)
	return ParseType(rdr)
}

func ParseType(rdr *bufio.Reader) interface{} {
	b, err := rdr.Peek(1)
	if err != nil {
		panic("failed to read byte")
	}

	switch b[0] {
	case 100: // d (map)
		return ParseMap(rdr)
	case 49, 50, 51, 52, 53, 54, 55, 56, 57: //number
		return ParseString(rdr)
	case 105: // i (number)
		return ParseNumber(rdr)
	case 108:
		return ParseList(rdr)
	}

	panic("Invalid bencoded string")
}

func ParseMap(rdr *bufio.Reader) map[string]interface{} {
	m := make(map[string]interface{})
	b := byte(0)

	// Advance reader past dict marker
	rdr.Discard(1)

	for b != 101 {
		k := ParseType(rdr).(string)
		m[k] = ParseType(rdr)
		ba, err := rdr.Peek(1)
		if err != nil {
			panic("Peek error in parsemap")
		}
		b = ba[0]
	}

	rdr.Discard(1)
	return m
}

func ParseList(rdr *bufio.Reader) []interface{} {
	l := []interface{}{}
	b := byte(0)

	// Advance reader past dict marker
	rdr.Discard(1)

	for b != 101 {
		l = append(l, ParseType(rdr))
		ba, err := rdr.Peek(1)
		if err != nil {
			panic("Peek error in parsemap")
		}
		b = ba[0]
	}

	rdr.Discard(1)
	return l
}

func ParseNumber(rdr *bufio.Reader) uint64 {
	rdr.Discard(1)
	cb, err := rdr.ReadString(byte(101))
	iv, err := strconv.ParseUint(strings.TrimSuffix(cb, "e"), 10, 64)

	if err != nil {
		panic("Failed to parse number")
	}

	return iv
}

func ParseLength(rdr *bufio.Reader) uint64 {
	cb, err := rdr.ReadString(byte(58))
	iv, err := strconv.ParseUint(strings.TrimSuffix(cb, ":"), 10, 64)

	if err != nil {
		panic("Failed to parse number")
	}

	return iv
}

func ParseString(rdr *bufio.Reader) string {
	l := ParseLength(rdr)
	s := []byte{}

	for i := uint64(0); i < l; i++ {
		b, err := rdr.ReadByte()
		if err != nil {
			panic("Failed to read byte")
		}
		s = append(s, b)
	}

	return string(s)
}
