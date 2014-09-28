package record

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

const ()

type Reader struct {
	ra io.ReaderAt
}

func NewReader(r io.ReaderAt) *Reader {
	return &Reader{ra: r}
}

func (r *Reader) ReadAt(index int64) (io.Reader, error) {
	//TODO change 4 to const
	crcSlice := make([]byte, 4, 4)
	valLenSlice := make([]byte, 4, 4)
	if _, err := r.ra.ReadAt(crcSlice, index); err != nil {
		return nil, err
	}

	if _, err := r.ra.ReadAt(valLenSlice, index+4); err != nil {
		return nil, err
	}

	crc := binary.LittleEndian.Uint32(crcSlice)
	valLen := binary.LittleEndian.Uint32(valLenSlice)

	value := make([]byte, valLen, valLen)
	if _, err := r.ra.ReadAt(value, index+4+4); err != nil {
		return nil, err
	}

	if crc != crc32.Update(crc32.Checksum(valLenSlice, crcTable), crcTable, value) {
		return nil, errors.New("crc unmatch, data corrupted")
	}

	return bytes.NewBuffer(value), nil
}
