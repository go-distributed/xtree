package record

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"io"
	"io/ioutil"
	"testing"
)

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
	valLenSlice := make([]byte, 4, 4)
	crcSlice := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(make([]byte, 4, 4), uint32(realLen))
	crc := crc32.Update(crc32.Checksum(valLenSlice, crcTable), crcTable, realValue)
	binary.LittleEndian.PutUint32(crcSlice, crc)

	buf.Write(make([]byte, 9, 9))
	buf.Write(crcSlice)
	buf.Write(valLenSlice)
	buf.Write(realValue)

	r := NewReader(bytes.NewReader(buf.Bytes()))
	valReader, err = r.ReadAt(9)
	if err != nil {
		t.Errorf("error")
	}

	readValue, err = ioutil.ReadAll(valReader)
	if err != nil {
		t.Errorf("error")
	}

	if ok := testEq(realValue, readValue); ok != true {
		t.Errorf("error")
	}

}
