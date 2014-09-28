package record

import "io"

type Writer struct {
}

func NewWriter(w io.WriteSeeker) *Writer {
	return nil
}

func (w *Writer) Append() (int64, io.Writer, error) {
	return 0, nil, nil
}

func (w *Writer) Flush() {

}
