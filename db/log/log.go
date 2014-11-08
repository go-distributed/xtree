package log

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/go-distributed/xtree/db/message"
)

const (
	logFilename = "records.log"
)

type DBLog struct {
	writeFile, readFile *os.File
	encoder             *encoder
}

func Create(dataDir string) (l *DBLog, err error) {
	l, err = newDBLog(path.Join(dataDir, logFilename), true)
	return
}

func newDBLog(logPath string, needCreate bool) (l *DBLog, err error) {
	var writeFile, readFile *os.File
	flag := os.O_WRONLY | os.O_APPEND | os.O_SYNC
	if needCreate {
		flag |= os.O_CREATE
	}
	if writeFile, err = os.OpenFile(logPath, flag, 0600); err != nil {
		return
	}
	if !needCreate {
		writeFile.Seek(0, os.SEEK_END)
	}
	if readFile, err = os.Open(logPath); err != nil {
		writeFile.Close()
		return
	}
	l = &DBLog{
		writeFile: writeFile,
		readFile:  readFile,
		encoder:   newEncoder(writeFile),
	}
	return
}

func Reuse(dataDir string,
	setLog func(*DBLog),
	replayRecord func(*message.Record) error) (err error) {
	var l *DBLog
	if l, err = newDBLog(path.Join(dataDir, logFilename),
		false); err != nil {
		return
	}
	setLog(l)
	decoder := newDecoder(l.readFile)
	// TODO: parallel?
	for {
		r := new(message.Record)
		if err = decoder.decode(r); err != nil {
			if err == io.EOF {
				return nil
			}
			return
		}
		ret, _ := l.readFile.Seek(0, os.SEEK_CUR)
		fmt.Printf("%v %#v\n", ret, r)
		if err = replayRecord(r); err != nil {
			return
		}
	}
}

func Exist(dataDir string) bool {
	p := path.Join(dataDir, logFilename)
	_, err := os.Stat(p)
	return err == nil
}

func (l *DBLog) GetRecord(offset int64) (r *message.Record, err error) {
	if _, err = l.readFile.Seek(offset, 0); err != nil {
		return
	}
	decoder := newDecoder(l.readFile)
	r = new(message.Record)
	err = decoder.decode(r)
	return
}

func (l *DBLog) Append(r *message.Record) (offset int64, err error) {
	if offset, err = l.writeFile.Seek(0, os.SEEK_CUR); err != nil {
		return
	}
	if err = l.encoder.encode(r); err != nil {
		return
	}
	err = l.encoder.flush()
	return offset, err
}

func (l *DBLog) Close() (err error) {
	if err = l.readFile.Close(); err != nil {
		return
	}
	return l.writeFile.Close()
}
