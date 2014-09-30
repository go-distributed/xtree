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

type Encoder interface {
	Encode(w io.Writer, r *Record) error
}

type Decoder interface {
	Decode(reader io.Reader, r *Record) error
}

type LittleEndianEncoder struct{}

type LittleEndianDecoder struct{}

func encodeLength(r *Record) []byte {
	lBuf := make([]byte, sizeOfLength)
	binary.LittleEndian.PutUint32(lBuf, uint32(len(r.data)))
	return lBuf
}

func calculateCRC(r *Record) uint32 {
	crc := crc32.Checksum(encodeLength(r), crcTable)
	crc = crc32.Update(crc, crcTable, r.data)
	return crc
}

func encodeCRC(r *Record) []byte {
	crcBuf := make([]byte, sizeOfCRC)
	binary.LittleEndian.PutUint32(crcBuf, calculateCRC(r))
	return crcBuf
}

func (*LittleEndianEncoder) Encode(w io.Writer, r *Record) error {
	// Write CRC
	if _, err := w.Write(encodeCRC(r)); err != nil {
		return err
	}
	// Write length
	if _, err := w.Write(encodeLength(r)); err != nil {
		return err
	}
	// Write data
	if _, err := w.Write(r.data); err != nil {
		return err
	}
	return nil
}

func (*LittleEndianDecoder) Decode(rd io.Reader, r *Record) error {
	var crc, length uint32

	// Read CRC
	err := binary.Read(rd, binary.LittleEndian, &crc)
	if err != nil {
		return err
	}
	// Read length
	err = binary.Read(rd, binary.LittleEndian, &length)
	if err != nil {
		return err
	}
	// Read data
	r.data = make([]byte, length)
	_, err = io.ReadFull(rd, r.data)
	if err != nil {
		return err
	}
	if crc != calculateCRC(r) {
		return errors.New("crc unmatch, data corrupted")
	}
	return nil
}
