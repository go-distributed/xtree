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

type RecordEncoder struct {
	w io.Writer
}

func NewRecordEncoder(w io.Writer) *RecordEncoder {
	return &RecordEncoder{w}
}

type RecordDecoder struct {
	r io.Reader
}

func NewRecordDecoder(r io.Reader) *RecordDecoder {
	return &RecordDecoder{r}
}

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

func (encoder *RecordEncoder) Encode(r *Record) error {
	// Write CRC
	if _, err := encoder.w.Write(encodeCRC(r)); err != nil {
		return err
	}
	// Write length
	if _, err := encoder.w.Write(encodeLength(r)); err != nil {
		return err
	}
	// Write data
	if _, err := encoder.w.Write(r.data); err != nil {
		return err
	}
	return nil
}

func (decoder *RecordDecoder) Decode(r *Record) error {
	var crc, length uint32

	// Read CRC
	err := binary.Read(decoder.r, binary.LittleEndian, &crc)
	if err != nil {
		return err
	}
	// Read length
	err = binary.Read(decoder.r, binary.LittleEndian, &length)
	if err != nil {
		return err
	}
	// Read data
	r.data = make([]byte, length)
	_, err = io.ReadFull(decoder.r, r.data)
	if err != nil {
		return err
	}
	if crc != calculateCRC(r) {
		return errors.New("crc unmatch, data corrupted")
	}
	return nil
}
