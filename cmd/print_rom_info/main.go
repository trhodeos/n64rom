package main

import (
	"encoding/binary"
	"fmt"
	"github.com/trhodeos/n64rom"
	"os"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	h, err := n64rom.ParseHeader(f, binary.BigEndian)
	fmt.Printf("%+v\n", h)
}
