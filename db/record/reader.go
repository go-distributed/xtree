package record

import "io"

type Reader struct {
}

func NewReader(r io.ReaderAt) *Reader {

}

func (r *Reader) ReadAt(index int64) (*io.Reader, error) {

}
