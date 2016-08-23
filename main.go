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
	case 100: // d (map)
		fmt.Println("Parse map")
		return ParseMap(rdr)
	case 49, 50, 51, 52, 53, 54, 55, 56, 57: //number
		fmt.Println("Parse string")
		l := ParseLength(rdr).IntValue()
		return ParseString(rdr, l)
	case 105: // i (number)
		return ParseNumber(rdr)
	}

	return Value{}
}

//TODO: only works for one key value pair... move lines 47-49 into this function
func ParseMap(rdr *bufio.Reader) Value {
	m := make(map[string]Value)

	// Advance reader past dict marker
	b, err := rdr.ReadByte()
	fmt.Println(string(b))
	if err != nil {
		panic("ReadByte error in ParseMap")
	}
	for b != 101 {
		k := ParseType(rdr)
		if k.ValueType() != "string" {
			panic("Map key must be a string.")
		}

		m[k.StringValue()] = ParseType(rdr)
		ba, err := rdr.Peek(1)
		fmt.Println("peek in parsemap: " + string(ba[0]))
		if err != nil {
			panic("Peek error in parsemap")
		}
		b = ba[0]
	}
	rdr.Discard(1)
	return Value{valueType: "map", mapValue: m}
}

func ParseNumber(rdr *bufio.Reader) Value {
	rdr.Discard(1)
	cb, err := rdr.ReadString(byte(101))
	fmt.Println("read string in parsenumber: " + string(cb))
	iv, err := strconv.ParseUint(strings.TrimSuffix(cb, "e"), 10, 64)

	if err != nil {
		panic("Failed to parse number")
	}

	return Value{valueType: "int", intValue: iv}
}

func ParseLength(rdr *bufio.Reader) Value {
	cb, err := rdr.ReadString(byte(58))
	fmt.Println(string(cb))
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
		fmt.Println(string(b))
		if err != nil {
			panic("Failed to read byte")
		}
		s = append(s, b)
	}

	return Value{valueType: "string", stringValue: string(s)}
}
