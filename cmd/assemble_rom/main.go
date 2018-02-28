package main

import (
	"github.com/trhodeos/n64rom"
	"os"
)

func main() {

	romheader, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	bootloader, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	font, err := os.Open(os.Args[3])
	if err != nil {
		panic(err)
	}
	writer, err := os.Create(os.Args[4])
	if err != nil {
		panic(err)
	}
	rom, err := n64rom.NewRomFile(romheader, bootloader, font, 0)
	if err != nil {
		panic(err)
	}
	rom.Save(writer)
}
