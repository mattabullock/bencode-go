package bencode

import (
	"bufio"
	"fmt"
	"os"
)

func test() {
	if len(os.Args) < 2 {
		panic("No arguments")
	}

	args := os.Args[1:]
	f, err := os.Open(args[0])
	if err != nil {
		panic("something broke")
	}

	reader := bufio.NewReader(f)
	enc, err := reader.ReadBytes('\n')
	if err != nil {
		panic("Failed to read from file")
	}

	m := Decode(enc)

	fmt.Println("Decoded: " + f.Name())
	fmt.Println(m)

	r := Encode(m)
	out, err := os.Create("temp")
	check(err)
	n2, err := out.Write(r)
	check(err)

	fmt.Printf("wrote %d bytes\n", n2)
	fmt.Println("Encoded to: " + out.Name())
	fmt.Println(r)
}
