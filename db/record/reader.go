package record

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

type Reader struct {
	ra io.ReaderAt
}

func NewReader(r io.ReaderAt) *Reader {
	return &Reader{ra: r}
}

func (r *Reader) ReadAt(index int64) (io.Reader, error) {
	crcSlice := make([]byte, sizeOfCRC, sizeOfCRC)
	valLenSlice := make([]byte, sizeOfLength, sizeOfLength)
	if _, err := r.ra.ReadAt(crcSlice, index); err != nil {
		return nil, err
	}

	if _, err := r.ra.ReadAt(valLenSlice, index+int64(sizeOfCRC)); err != nil {
		return nil, err
	}

	crc := binary.LittleEndian.Uint32(crcSlice)
	valLen := binary.LittleEndian.Uint32(valLenSlice)

	value := make([]byte, valLen, valLen)
	if _, err := r.ra.ReadAt(value, index+int64(sizeOfCRC+sizeOfLength)); err != nil {
		return nil, err
	}

	if crc != crc32.Update(crc32.Checksum(valLenSlice, crcTable), crcTable, value) {
		return nil, errors.New("crc unmatch, data corrupted")
	}

	return bytes.NewBuffer(value), nil
}
