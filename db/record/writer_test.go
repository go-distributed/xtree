package record

import "testing"
import "bytes"

//import "fmt"

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
	w, err := NewWriter(m)
	if err != nil {
		t.Fatalf("Cannot create writer")
	}
	offset, value, _ := w.Append()
	if offset != 0 {
		t.Fatalf("Initial offset is not 0")
	}
	s := "test string"
	value.Write([]byte(s))
	offset, value, err = w.Append()
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
	if m.Len() != 8+len(s) {
		t.Fatalf("File length not correct")
	}
	if len(m.Bytes()[4:]) != 4+len(s) {
		t.Fatalf("%d != %d", len(m.Bytes()[4:]), 4+len(s))
	}
}
