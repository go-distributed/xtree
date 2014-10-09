package recordio

import "io"

type Appender interface {
	Append(Encoder) (int64, error)
}

type appender struct {
	w io.WriteSeeker
}

func NewAppender(w io.WriteSeeker) Appender {
	return &appender{w}
}

// Not thread-safe
func (ap *appender) Append(r Encoder) (offset int64, err error) {
	offset, err = ap.w.Seek(0, 2)
	if err != nil {
		return -1, err
	}

	err = r.EncodeTo(ap.w)
	if err != nil {
		return -1, err
	}

	return offset, nil
}
