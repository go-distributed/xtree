package record

import "io"

type Reader struct {
}

func NewReader(r io.ReaderAt) *Reader {
	return nil
}

func (r *Reader) ReadAt(offset int64) (io.Reader, error) {
	return nil, nil
}
