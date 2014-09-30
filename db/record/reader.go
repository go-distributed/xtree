package record

import (
	"io"
)

type Reader struct {
	rs io.ReadSeeker
}

func NewReader(rs io.ReadSeeker) *Reader {
	return &Reader{rs: rs}
}

func (rd *Reader) ReadAt(index int64) (io.Reader, error) {
	_, err := rd.rs.Seek(index, 0)
	if err != nil {
		return nil, err
	}
	return rd.rs, nil
}
