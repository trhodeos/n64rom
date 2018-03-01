package main

import (
	"fmt"
	"github.com/trhodeos/n64rom"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <bootloader> <font> <outpath>\n", os.Args[0])
		return
	}
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

	blankData := make([]byte, 0x400000)

	rom.WriteAt(blankData, 0x1000)

	total, err := rom.Save(writer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Saved %d bytes to %s.\n", total, os.Args[3])
}
