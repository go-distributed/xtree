package log

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/go-distributed/xtree/db/message"
)

type encoder struct {
	bw *bufio.Writer
}

func newEncoder(w io.Writer) *encoder {
	return &encoder{bufio.NewWriter(w)}
}

func (e *encoder) encode(r *message.Record) (err error) {
	var data []byte
	if data, err = r.Marshal(); err != nil {
		return
	}
	if err = writeInt64(e.bw, int64(len(data))); err != nil {
		return
	}
	_, err = e.bw.Write(data)
	return
}

func (e *encoder) flush() error {
	return e.bw.Flush()
}

func writeInt64(w io.Writer, n int64) error {
	return binary.Write(w, binary.LittleEndian, n)
}
