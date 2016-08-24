package main

import (
	"bufio"
	"fmt"
	"os"
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

	reader := bufio.NewReader(f)
	enc, err := reader.ReadBytes('\n')
	if err != nil {
		panic("Failed to read from file")
	}

	r := ParseBencode(enc)
	fmt.Println(r)
}
