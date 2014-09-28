package record

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
	"io/ioutil"
	"testing"
)

const randOffset = 9

func TestSimpleRead(t *testing.T) {

	testEq := func(a, b []byte) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}

	var buf bytes.Buffer
	var err error
	var valReader io.Reader
	var readValue []byte
	realValue := []byte("Hello, World")
	realLen := len(realValue)
	crcSlice := make([]byte, sizeOfCRC, sizeOfCRC)
	realLenSlice := make([]byte, sizeOfLength, sizeOfLength)
	binary.LittleEndian.PutUint32(realLenSlice, uint32(realLen))
	crc := crc32.Update(crc32.Checksum(realLenSlice, crcTable), crcTable, realValue)
	binary.LittleEndian.PutUint32(crcSlice, crc)

	buf.Write(make([]byte, randOffset, randOffset))
	buf.Write(crcSlice)
	buf.Write(realLenSlice)
	buf.Write(realValue)

	r := NewReader(bytes.NewReader(buf.Bytes()))
	valReader, err = r.ReadAt(randOffset)
	if err != nil {
		t.Fatalf("error on new reader: %s", err.Error())
	}

	readValue, err = ioutil.ReadAll(valReader)
	if err != nil {
		t.Fatalf("error on read all")
	}

	if ok := testEq(realValue, readValue); ok != true {
		t.Fatalf("value not equal")
	}

}
