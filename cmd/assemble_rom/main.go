package main

import (
	"fmt"
	"github.com/trhodeos/n64rom"
	"os"
)

func main() {
	bootloader, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	font, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	writer, err := os.Create(os.Args[3])
	if err != nil {
		panic(err)
	}

	rom, err := n64rom.NewBlankRomFile(bootloader, font, 0)
	if err != nil {
		panic(err)
	}

	total, err := rom.Save(writer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Saved %d bytes to %s.", total, os.Args[3])
}
