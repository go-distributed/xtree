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
	if err := l.f.Close(); err != nil {
		return err
	}
	return os.Remove(l.f.Name())
}

func (l *Log) GetRecord(offset int64) (r *message.Record, err error) {
	if _, err = l.f.Seek(offset, 0); err != nil {
		return
	}
	decoder := newDecoder(l.f)
	r = &message.Record{}
	err = decoder.decode(r)
	return
}

func (l *Log) Append(r *message.Record) (offset int64, err error) {
	if offset, err = l.f.Seek(0, 2); err != nil {
		return
	}
	encoder := newEncoder(l.f)
	if err = encoder.encode(r); err != nil {
		return
	}
	err = encoder.flush()
	return offset, err
}
