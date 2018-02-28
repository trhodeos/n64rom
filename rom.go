package n64rom

import (
	"encoding/binary"
	"io"
)

type TvType int

const (
	Unknown TvType = iota
	Ntsc
	Pal
)

type Country struct {
	Name   string
	TvType TvType
}

// From http://en64.shoutwiki.com/wiki/ROM
var CountryInfo = map[uint8]Country{
	0x0:  {"0", Ntsc},
	0x37: {"Beta", Ntsc},
	0x41: {"Asian (NTSC)", Ntsc},
	0x42: {"Brazilian", Ntsc},
	0x43: {"Chinese", Ntsc},
	0x44: {"German", Pal},
	0x45: {"North America", Ntsc},
	0x46: {"French", Pal},
	0x47: {"Gateway 64", Ntsc},
	0x48: {"Dutch", Pal},
	0x49: {"Italian", Pal},
	0x4A: {"Japanese", Ntsc},
	0x4B: {"Korean", Ntsc},
	0x4C: {"Gateway 64", Pal},
	0x4E: {"Canadian", Ntsc},
	0x50: {"European", Pal},
	0x53: {"Spanish", Pal},
	0x55: {"Australian", Pal},
	0x57: {"Scandinavian", Pal},
	0x58: {"European", Pal},
	0x59: {"European", Pal},
}

type Header struct {
	X1 uint8
	X2 uint8
	X3 uint8
	X4 uint8

	ClockRate      uint32
	BootAddress    uint32
	Release        uint32
	Crc1           uint32
	Crc2           uint32
	Unknown0       uint64
	Name           [20]uint8
	Unknown2       uint32
	RomType        uint8
	GameId         uint16
	RegionLanguage uint8
	CartId         uint16
	CountryCode    uint8
	Version        uint8
}

func ParseHeader(r io.Reader, bo binary.ByteOrder) (Header, error) {
	out := Header{}
	err := binary.Read(r, bo, &out)
	return out, err
}
