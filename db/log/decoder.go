package log

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/go-distributed/xtree/db/message"
)

type decoder struct {
	br *bufio.Reader
}

func newDecoder(r io.Reader) *decoder {
	return &decoder{bufio.NewReader(r)}
}

func (d *decoder) decode(r *message.Record) (err error) {
	var l int64
	if l, err = readInt64(d.br); err != nil {
		return
	}
	data := make([]byte, l)
	if _, err = io.ReadFull(d.br, data); err != nil {
		return
	}
	return r.Unmarshal(data)
}

func readInt64(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.LittleEndian, &n)
	return
}
