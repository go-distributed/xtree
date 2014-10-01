package record

import (
	"bytes"
	"io"
)

type Reader struct {
	rs  io.ReadSeeker
	dec Decoder
}

func NewReader(rs io.ReadSeeker, dec Decoder) *Reader {
	return &Reader{rs, dec}
}

func (rd *Reader) ReadFromIndex(index int64) (io.Reader, error) {
	_, err := rd.rs.Seek(index, 0)
	if err != nil {
		return nil, err
	}

	rec := Record{}
	err = rd.dec.Decode(rd.rs, &rec)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(rec.data), nil
}
