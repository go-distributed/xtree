package log

import (
	"io/ioutil"
	"os"

	"github.com/go-distributed/xtree/db/message"
)

type Log struct {
	f *os.File
}

func Create() (*Log, error) {
	f, err := ioutil.TempFile("", "backend")
	if err != nil {
		return nil, err
	}

	return &Log{f}, nil
}

func (l *Log) Destroy() error {
	return os.Remove(l.f.Name())
}

func (l *Log) GetRecord(offset int64) (*message.Record, error) {
	_, err := l.f.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	decoder := newDecoder(l.f)
	r := &message.Record{}
	decoder.decode(r)

	return r, nil
}

func (l *Log) Append(r *message.Record) (offset int64, err error) {
	offset, err = l.f.Seek(0, 2)
	if err != nil {
		return -1, err
	}

	encoder := newEncoder(l.f)

	err = encoder.encode(r)
	encoder.flush()
	return offset, err
}
