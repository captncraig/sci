package resources

import (
	"bytes"
	"compress/lzw"
	"fmt"
	"io/ioutil"

	"github.com/captncraig/sci"
)

type ResourceType byte

const (
	View ResourceType = iota
	Picture
	Script
	Text
	Sound
	Vocab
	Font
	Cursor
	Patch
)

func (rt ResourceType) String() string {
	switch rt {
	case View:
		return "view"
	case Picture:
		return "picture"
	case Script:
		return "script"
	case Text:
		return "text"
	case Sound:
		return "sound"
	case Vocab:
		return "vocab"
	case Font:
		return "font"
	case Cursor:
		return "cursor"
	case Patch:
		return "patch"
	}
	return "?"
}

type ResourceHeader struct {
	Type       ResourceType
	ID         uint16
	FileNumber byte
	Offset     uint32

	LoadError        string
	CompressedSize   uint16
	DecompressedSize uint16
	Method           CompressionMethod
	Data             []byte
}

func (r *ResourceHeader) String() string {
	return fmt.Sprintf("%s.%03d, resource.%03d @ $%08x", r.Type, r.ID, r.FileNumber, r.Offset)
}

type CompressionMethod uint16

const (
	Uncompressed CompressionMethod = iota
	LZW
	Huffman
)

func (cm CompressionMethod) String() string {
	switch cm {
	case Uncompressed:
		return "none"
	case LZW:
		return "lzw"
	case Huffman:
		return "huffman"
	}
	return "?"
}

//The SCI0 map file format is pretty simple: It consists of 6-byte entries, terminated by the sequence 0xffff ffff ffff.
// The first 2 bytes, interpreted as little endian 16 bit integer, encode resource type (high 5 bits) and number (low 11 bits).
//The next 4 bytes are a 32 bit LE integer that contains the resource file number in the high 6 bits,
//and the absolute offset within the file in the low 26 bits.
//SCI0 performs a linear search to find the resource; however, multiple entries may match the search, since resources may be present
//more than once (the inverse mapping is not injective).

// ReadMap will parse the resource map and return a list of resource pointers. If a non-nil loader is given,
// it will also read the associated data, decompress it, and store it in the header.
func ReadMap(dat []byte, loader sci.Loader) ([]*ResourceHeader, error) {
	// make sure it ends with 0xff * 6 for SCI0
	if len(dat)%6 != 0 || len(dat) == 0 {
		return nil, fmt.Errorf("resource map should be a multiple of 6 bytes long")
	}
	if string(dat[len(dat)-6:]) != "\xff\xff\xff\xff\xff\xff" {
		return nil, fmt.Errorf("resource map should end with 0xffff ffff ffff")
	}
	dat = dat[:len(dat)-6]
	rs := make([]*ResourceHeader, 0, len(dat)/6)
	seenKeys := map[uint16]bool{}
	for i := 0; i < len(dat)-5; i += 6 {
		key := uint16(dat[i]) | uint16(dat[i+1])<<8
		if seenKeys[key] {
			continue
		}
		seenKeys[key] = true
		rec := &ResourceHeader{
			Type:       ResourceType(dat[i+1] >> 3),
			ID:         read16(dat, i) & 0x07ff,
			FileNumber: dat[i+5] >> 2,
			Offset:     read32(dat, i+2) & 0x03ffffff,
		}
		rs = append(rs, rec)
		if loader != nil {
			rec.Load(loader)
		}

	}
	return rs, nil
}

func read16(dat []byte, i int) uint16 {
	return uint16(dat[i]) | uint16(dat[i+1])<<8
}
func read32(dat []byte, i int) uint32 {
	return uint32(dat[i]) | uint32(dat[i+1])<<8 | uint32(dat[i+2])<<16 | uint32(dat[i+3])<<24
}

func (r *ResourceHeader) Load(l sci.Loader) {
	defer func() {
		if r.LoadError != "" {
			fmt.Println(r.LoadError)
		}
	}()
	dat, err := l.GetFile(fmt.Sprintf("RESOURCE.%03d", r.FileNumber))
	if err != nil {
		r.LoadError = err.Error()
		return
	}
	dat = dat[r.Offset:]
	if len(dat) < 8 {
		r.LoadError = fmt.Sprintf("Not enough data in resource.%03d to get offset 0x%x", r.FileNumber, r.Offset)
		return
	}
	id := read16(dat, 0)
	if id != r.ID|uint16(r.Type)<<11 {
		r.LoadError = fmt.Sprintf("ID does not match resource data at offset 0x%x", r.Offset)
		return
	}
	r.CompressedSize = read16(dat, 2)
	r.CompressedSize -= 4 // decompressed size and method should not count here
	r.DecompressedSize = read16(dat, 4)
	r.Method = CompressionMethod(read16(dat, 6))
	dat = dat[8:]
	if len(dat) < int(r.CompressedSize) {
		r.LoadError = fmt.Sprintf("Not enough data in resource.%03d to read %d compressed bytes at 0x%x", r.FileNumber, r.CompressedSize, r.Offset+8)
		return
	}
	dat = dat[:r.CompressedSize]
	switch r.Method {
	case Uncompressed:
		break
	case LZW:
		rd := lzw.NewReader(bytes.NewReader(dat), lzw.LSB, 8)
		dat, err = ioutil.ReadAll(rd)
		if err != nil {
			r.LoadError = fmt.Sprintf("Errror in lzw decompression: %s", err)
			return
		}
	case Huffman:
		fmt.Println(r)
		dat = huffmanDecode(dat, r.DecompressedSize)
	default:
		r.LoadError = fmt.Sprintf("Unimplemented decompression: %s", r.Method)
		return
	}
	r.Data = dat
	if len(r.Data) != int(r.DecompressedSize) {
		r.LoadError = fmt.Sprintf("Data length (%d) does not match expected decompressed size (%d)", len(r.Data), r.DecompressedSize)
		return
	}
}
