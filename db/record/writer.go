package record

import (
	"io"
	"io/ioutil"
)

type flusher interface {
	Flush() error
}

const bufferSize = 32 * 1024

type Writer struct {
	// Underlying writer
	w   io.WriteSeeker
	enc Encoder
}

func NewWriter(w io.WriteSeeker, enc Encoder) *Writer {
	return &Writer{w, enc}
}

// Not thread-safe
func (wt *Writer) Append(d io.Reader) (offset int64, err error) {
	var data []byte

	offset, err = wt.w.Seek(0, 2)
	if err != nil {
		return 0, err
	}
	data, err = ioutil.ReadAll(d)
	rec := Record{data: data}
	err = wt.enc.Encode(wt.w, &rec)
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func (wt *Writer) Flush() error {
	f, ok := wt.w.(flusher)
	if ok {
		return f.Flush()
	}
	return nil
}
