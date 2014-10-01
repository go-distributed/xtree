package record

import "testing"
import "bytes"
import "errors"
import "io"
import "io/ioutil"

type memFile struct {
	b   *bytes.Buffer
	off int64
}

var (
	testString1 = "testString1"
	testString2 = "testString2"
)

func (m *memFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		if offset >= int64(m.b.Len()) {
			return 0, errors.New("invalid offset")
		}
		m.off = offset
		return m.off, nil
	case 2:
		return int64(m.b.Len()), nil
	}
	return 0, errors.New("unimplemented")
}

func (m *memFile) Write(b []byte) (n int, err error) {
	return m.b.Write(b)
}

func (m *memFile) Read(b []byte) (n int, err error) {
	n, err = bytes.NewBuffer(m.b.Bytes()[m.off:]).Read(b)
	if err == nil {
		m.off += int64(n)
	}
	return
}

type simpleEncoder struct{}

func (*simpleEncoder) Encode(w io.Writer, r *Record) error {
	_, err := w.Write(r.data)
	return err
}

type simpleDecoder struct{}

func (*simpleDecoder) Decode(rd io.Reader, r *Record) (err error) {
	r.data, err = ioutil.ReadAll(rd)
	return
}

func TestWriter(t *testing.T) {
	m := &memFile{new(bytes.Buffer), 0}
	w := NewWriter(m, new(simpleEncoder))
	offset, err := w.Append(bytes.NewBufferString(testString1))
	if err != nil {
		t.Fatalf("cannot append: %s", err)
	}
	if offset != 0 {
		t.Fatalf("initial Offset is not zero")
	}
	offset, err = w.Append(bytes.NewBufferString(testString2))
	if !bytes.Equal(m.b.Bytes(), []byte(testString1+testString2)) {
		t.Fatalf("got unpected results")
	}
}

func TestReader(t *testing.T) {
	m := &memFile{bytes.NewBufferString(testString1), 0}
	r := NewReader(m, new(simpleDecoder))
	reader, err := r.ReadFromIndex(0)
	if err != nil {
		t.Fatalf("cannot append: %s", err)
	}
	results, err := ioutil.ReadAll(reader)
	if !bytes.Equal(results, []byte(testString1)) {
		t.Fatalf("got unpected results")
	}
}
