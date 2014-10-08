package recordio

import (
	"encoding/binary"
	"io"
)

const (
	sizeOfLength = 4
)

type Record struct {
	Data []byte
}

func (r *Record) encodeTo(wr io.Writer) error {
	// Write length
	lBuf := make([]byte, sizeOfLength)
	binary.LittleEndian.PutUint32(lBuf, uint32(len(r.Data)))
	if _, err := wr.Write(lBuf); err != nil {
		return err
	}
	// Write data
	if _, err := wr.Write(r.Data); err != nil {
		return err
	}

	return nil
}

func (r *Record) decodeFrom(rd io.Reader) error {
	var length uint32
	// Read length
	err := binary.Read(rd, binary.LittleEndian, &length)
	if err != nil {
		return err
	}
	// Read data
	r.Data = make([]byte, length)
	_, err = io.ReadFull(rd, r.Data)
	if err != nil {
		return err
	}

	return nil
}
