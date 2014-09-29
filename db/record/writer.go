package record

import "io"

type flusher interface {
	Flush() error
}

const bufferSize = 32 * 1024

type Writer struct {
	// Underlying writer
	w io.WriteSeeker
}

func NewWriter(w io.WriteSeeker) *Writer {
	return &Writer{w}
}

// Not thread-safe
func (w *Writer) Append() (offset int64, value io.Writer, err error) {
	offset, err = w.w.Seek(0, 2)
	value = w.w
	err = nil
	return
}

func (w *Writer) Flush() error {
	f, ok := w.w.(flusher)
	if ok {
		return f.Flush()
	}
	return nil
}
