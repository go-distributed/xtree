package record

import "testing"
import "bytes"
import "encoding/binary"

type memFile struct {
	bytes.Buffer
}

func (m *memFile) Seek(offset int64, whence int) (int64, error) {
	return int64(m.Len()), nil
}

func TestWriter(t *testing.T) {
	var b bytes.Buffer
	m := &memFile{b}
	if m.Len() != 0 {
		t.Fatalf("Initial file length is not zero")
	}
	w := NewWriter(m)
	offset, value, _ := w.Append()
	if offset != 0 {
		t.Fatalf("Initial offset is not 0")
	}
	s := "test string"
	NewRecordEncoder(value).Encode(&Record{[]byte(s)})
	offset, value, err := w.Append()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if offset == 0 {
		t.Fatalf("Offset is still 0 after write")
	}
	err = w.Flush()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if m.Len() != sizeOfLength+sizeOfCRC+len(s) {
		t.Fatalf("File length is incorrect")
	}
	if len(m.Bytes()[sizeOfCRC:]) != sizeOfCRC+len(s) {
		t.Fatalf("Bytes length is incorrect: expect: %d, got: %d", len(m.Bytes()[sizeOfCRC:]), sizeOfCRC+len(s))
	}
	if int(binary.LittleEndian.Uint32(m.Bytes()[sizeOfCRC:sizeOfCRC+sizeOfLength+1])) != len(s) {
		t.Fatalf("Data length is incorrect")
	}
}
