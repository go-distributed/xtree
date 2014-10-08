package recordio

import (
	"encoding/binary"
	"io"
)

const (
	sizeOfLength = 4
)

// Encoder is an interface implemented by an object that can
// encode itself into a binary representation with an io.Writer
type Encoder interface {
	EncodeTo(wr io.Writer) error
}

// Decoder is an interface implemented by an object that can
// decodes itself from a binary representation provided by
// an io.Reader
type Decoder interface {
	DecodeFrom(rd io.Reader) error
}

// Record is a struct that holds some binary data and can be
// encoded into/decoded from the following binary form
// [lenth of record]                  : 4 bytes
// [content of record]
type Record struct {
	Data []byte
}

// EncodeTo encodes a binary data in record object to an io.Writer
func (r *Record) EncodeTo(wr io.Writer) error {
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

// DecodeFrom decodes binary data out of io.Reader
func (r *Record) DecodeFrom(rd io.Reader) error {
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
