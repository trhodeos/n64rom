package n64rom

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

type RomFile struct {
  header Header
	bootloaderAndFont []byte
	objects           map[int64][]byte
	fillValue         byte
}

const maxHeaderSize = 0x40
const bootloaderStart = 0x40
const maxBootloaderSize = codeStart - maxHeaderSize
const codeStart = 0x1000

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
	err = checkData("Bootloader+Font", &r.bootloaderAndFont, maxBootloaderSize)
	if err != nil {
		return err
	}
	return nil
}

func (r RomFile) fillAt(o io.WriterAt, start int, end int) (int, error) {
	fill := bytes.Repeat([]byte{r.fillValue}, end-start)
	return o.WriteAt(fill, int64(start))
}

func (r RomFile) getHeaderBytes() ([]byte, error) {
        buffer := bytes.Buffer.NewBuffer([]byte)
	err := binary.Write(buffer, ByteOrder.BigEndian, r.header)
        return buffer.Bytes(), err
}

func NewBlankRomFile(bootloader io.Reader, font io.Reader, fill byte) (RomFile, error) {
  // TODO: search PATH for bootloader and font files
  return NewRomFile(GetBlankHeader(), bootloader, font, fill)
}

func NewRomFile(header Header, bootloader io.Reader, font io.Reader, fill byte) (RomFile, error) {
	bootloaderBytes, err := ioutil.ReadAll(bootloader)
	if err != nil {
		return RomFile{}, err
	}
	fontBytes, err := ioutil.ReadAll(font)
	if err != nil {
		return RomFile{}, err
	}

	rom := RomFile{header: header, bootloaderAndFont: append(bootloaderBytes, fontBytes...), fillValue: fill}
	err = rom.checkValidity()
	return rom, err
}

func (r *RomFile) WriteAt(p []byte, i int64) error {
	if i < codeStart {
		return errors.New(
			fmt.Sprintf("Cannot write at %d: This would overwrite bootloader before %d", i, codeStart))
	}
	r.objects[i] = p
	return nil
}

func (r *RomFile) Save(o io.WriterAt) (int, error) {
	total := 0
	n, err := o.WriteAt(r.header, 0x0)
	total += n
	if err != nil {
		return total, err
	}
	n, err = r.fillAt(o, len(r.header), maxHeaderSize)
	total += n
	if err != nil {
		return total, err
	}
	n, err = o.WriteAt(r.bootloaderAndFont, bootloaderStart)
	total += n
	if err != nil {
		return total, err
	}
	n, err = r.fillAt(o, maxHeaderSize+len(r.bootloaderAndFont), codeStart)
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

	// TODO update checksum
	return total, err
}
