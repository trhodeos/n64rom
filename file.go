package n64rom

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/trhodeos/ecoff"
	"io"
	"io/ioutil"
)

type RomFile struct {
	header            Header
	bootloaderAndFont []byte
	objects           map[int64][]byte
	fillValue         byte
}

const maxHeaderSize = 0x40
const bootloaderStart = 0x40
const maxBootloaderSize = CodeStart - maxHeaderSize
const CodeStart = 0x1000

func checkData(name string, data *[]byte, maxSize int) error {
	if data == nil {
		return errors.New(fmt.Sprintf("%s must be defined!", name))
	}
	if len(*data) == 0 {
		return errors.New(fmt.Sprintf("%s must have size greater than 0", name))
	}
	if maxSize > 0 && len(*data) > maxSize {
		return errors.New(fmt.Sprintf("%s must be less than or equal to 0x%x bytes!", name, maxSize))
	}
	return nil
}

func (r RomFile) checkValidity() error {
	//err := checkData("Bootloader+Font", &r.bootloaderAndFont, maxBootloaderSize)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (r RomFile) fillAt(o io.WriterAt, start int, end int) (int, error) {
	fill := bytes.Repeat([]byte{r.fillValue}, end-start)
	return o.WriteAt(fill, int64(start))
}

func (r RomFile) getHeaderBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, r.header)
	return buf.Bytes(), err
}

func NewBlankRomFile(fill byte) (RomFile, error) {
	// TODO: search PATH for bootloader and font files
	return NewRomFile(GetBlankHeader(), nil, nil, fill)
}

func getPreambleBytes(bootloader io.Reader, font io.Reader) ([]byte, error) {
	bootloaderBytes, err := ioutil.ReadAll(bootloader)
	if err != nil {
		return nil, err
	}
	bootloaderHeader, err := ecoff.ParseHeader(
		bytes.NewReader(bootloaderBytes), binary.BigEndian)
	if err != nil {
		return nil, err
	}
	if len(bootloaderHeader.SectionHeaders) != 2 {
		return nil, errors.New(fmt.Sprintf("Expected 2 sections in bootloader found %d.", len(bootloaderHeader.SectionHeaders)))
	}
	var sectionHeader *ecoff.SectionHeader
	for _, header := range bootloaderHeader.SectionHeaders {
		if string(header.Name[0:5]) == ".text" {
			sectionHeader = &header
			break
		}
	}
	if sectionHeader == nil {
		return nil, errors.New("Could not find section named '.text'.")
	}
	textSize := int64(sectionHeader.Size)
	bootloaderTextBytes := make([]byte, textSize)
	_, err = bytes.NewReader(bootloaderBytes).ReadAt(bootloaderTextBytes, textSize)
	if err != nil {
		return nil, err
	}
	fontBytes, err := ioutil.ReadAll(font)
	if err != nil {
		return nil, err
	}
	return append(bootloaderTextBytes, fontBytes...), nil
}

func NewRomFile(header Header, bootloader io.Reader, font io.Reader, fill byte) (RomFile, error) {
	var preamble []byte
	var err error
	if bootloader != nil && font != nil {
		preamble, err = getPreambleBytes(bootloader, font)
		if err != nil {
			return RomFile{}, err
		}
	} else {
		preamble = []byte{}
	}
	rom := RomFile{header: header, objects: map[int64][]byte{}, bootloaderAndFont: preamble, fillValue: fill}
	err = rom.checkValidity()
	return rom, err
}

func (r *RomFile) WriteAt(p []byte, i int64) error {
	if i < CodeStart {
		return errors.New(
			fmt.Sprintf("Cannot write at %d: This would overwrite bootloader before %d", i, CodeStart))
	}
	r.objects[i] = p
	return nil
}

func (r *RomFile) Save(o io.WriterAt) (int, error) {
	total := 0
	headerBytes, err := r.getHeaderBytes()
	if err != nil {
		return total, err
	}

	n, err := o.WriteAt(headerBytes, 0x0)
	total += n
	if err != nil {
		return total, err
	}
	n, err = r.fillAt(o, len(headerBytes), maxHeaderSize)
	total += n
	if err != nil {
		return total, err
	}
	n, err = o.WriteAt(r.bootloaderAndFont, bootloaderStart)
	total += n
	if err != nil {
		return total, err
	}
	n, err = r.fillAt(o, maxHeaderSize+len(r.bootloaderAndFont), CodeStart)
	total += n
	if err != nil {
		return total, err
	}
	for addr, obj := range r.objects {
		// TODO implement checks for overlaps
		n, err := o.WriteAt(obj, addr)
		total += n
		if err != nil {
			return total, err
		}
	}

	return total, err
}
