package record

import "io"

type Writer struct {
}

func NewWriter(w io.Writer) *Writer {

}

func (w *Writer) Append() (int64, *io.Writer, error) {

}

func (w *Writer) Flush() {

}
