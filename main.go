package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		panic("No arguments")
	}

	args := os.Args[1:]
	f, err := os.Open(args[0])
	if err != nil {
		panic("something broke")
	}

	Parse(f)
}

func Parse(f *os.File) {
	reader := bufio.NewReader(f)
	m := ParseType(reader)
	for key, value := range m.(map[string]interface{}) {
		fmt.Println("\""+key+"\", ", value)
	}
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

	return ""
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
