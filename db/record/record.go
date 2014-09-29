package record

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"unsafe"
)

const (
	sizeOfLength = int(unsafe.Sizeof(uint32(0)))
	sizeOfCRC    = int(unsafe.Sizeof(uint32(0)))
)

var crcTable = crc32.MakeTable(crc32.Koopman)

type Record struct {
	data []byte
}

type encoder interface {
	Encode(r Record) error
}

type RecordEncoder struct {
	w    io.Writer
	used bool
}

func NewRecordEncoder(w io.Writer) *RecordEncoder {
	return &RecordEncoder{w, false}
}

func (encoder *RecordEncoder) encodeLength(r *Record) []byte {
	lBuf := make([]byte, sizeOfLength)
	binary.LittleEndian.PutUint32(lBuf, uint32(len(r.data)))
	return lBuf
}

func (encoder *RecordEncoder) crc(r *Record) uint32 {
	crc := crc32.Checksum(encoder.encodeLength(r), crcTable)
	crc = crc32.Update(crc, crcTable, r.data)
	return crc
}

func (encoder *RecordEncoder) encodeCRC(r *Record) []byte {
	crcBuf := make([]byte, sizeOfCRC)
	binary.LittleEndian.PutUint32(crcBuf, encoder.crc(r))
	return crcBuf
}

func (encoder *RecordEncoder) Encode(r *Record) error {
	if encoder.used {
		return errors.New("encoder already used")
	}
	// Write CRC
	if _, err := encoder.w.Write(encoder.encodeCRC(r)); err != nil {
		return err
	}
	// Write length
	if _, err := encoder.w.Write(encoder.encodeLength(r)); err != nil {
		return err
	}
	// Write data
	if _, err := encoder.w.Write(r.data); err != nil {
		return err
	}
	encoder.used = true
	return nil
}
