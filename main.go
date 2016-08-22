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
	for key, value := range m.MapValue() {
		fmt.Println("Key:", key, "Value:", value)
	}
}

func ParseType(rdr *bufio.Reader) Value {
	b, err := rdr.Peek(1)
	if err != nil {
		panic("failed to read byte")
	}

	switch b[0] {
	case 100:
		return ParseMap(rdr)
	case 49, 50, 51, 52, 53, 54, 55, 56, 57:
		l := ParseNumber(rdr).IntValue()
		return ParseString(rdr, l)
	case 105:
		return ParseNumber(rdr)
	case 101:
		rdr.Discard(1)
		ParseType(rdr)
	}

	return Value{}
}

func ParseMap(rdr *bufio.Reader) Value {
	m := make(map[string]Value)

	// Advance reader past dict marker
	rdr.Discard(1)
	k := ParseType(rdr)
	if k.ValueType() != "string" {
		panic("Map key must be a string.")
	}

	m[k.StringValue()] = ParseType(rdr)
	return Value{valueType: "map", mapValue: m}
}

func ParseNumber(rdr *bufio.Reader) Value {
	cb, err := rdr.ReadString(byte(58))
	iv, err := strconv.ParseUint(strings.TrimSuffix(cb, ":"), 10, 64)

	if err != nil {
		panic("Failed to parse number")
	}

	return Value{valueType: "int", intValue: iv}
}

func ParseString(rdr *bufio.Reader, l uint64) Value {
	s := []byte{}
	for i := uint64(0); i < l; i++ {
		b, err := rdr.ReadByte()
		if err != nil {
			panic("Failed to read byte")
		}
		s = append(s, b)
	}

	return Value{valueType: "string", stringValue: string(s)}
}
