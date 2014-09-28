package record

import "io"
import "unsafe"
import "errors"
import "hash/crc32"
import "encoding/binary"

var crcTable = crc32.MakeTable(crc32.Koopman)

type flusher interface {
	Flush() error
}

const bufferSize = 32 * 1024

const (
	sizeOfLength = int(unsafe.Sizeof(uint32(0)))
	sizeOfCRC    = int(unsafe.Sizeof(uint32(0)))
)

type recordWriter struct {
	w      *Writer
	offset int64
	buf    []byte
}

func (r *recordWriter) Write(p []byte) (int, error) {
	if r.w.offset != r.offset {
		return 0, errors.New("stale writer")
	}
	r.buf = append(r.buf, p...)
	return len(p), nil
}

func (r *recordWriter) finalize() error {
	if len(r.buf) == 0 {
		return nil
	}
	l := len(r.buf)
	lBuf := make([]byte, sizeOfLength)
	binary.LittleEndian.PutUint32(lBuf, uint32(l))
	crc := crc32.Checksum(lBuf, crcTable)
	crc = crc32.Update(crc, crcTable, r.buf)
	crcBuf := make([]byte, sizeOfCRC)
	binary.LittleEndian.PutUint32(crcBuf, crc)
	// Write CRC
	if _, err := r.w.w.Write(crcBuf); err != nil {
		return err
	}
	// Write length
	if _, err := r.w.w.Write(lBuf); err != nil {
		return err
	}
	if _, err := r.w.w.Write(r.buf); err != nil {
		return err
	}
	newOffset, err := r.w.w.Seek(0, 2)
	if newOffset == 0 {
		panic("Cannot update offset")
	}
	r.w.offset = newOffset
	return err
}

type Writer struct {
	offset int64
	// Writer given out to the client to write data
	ongoingWrite *recordWriter
	// Underlying writer
	w io.WriteSeeker
}

func NewWriter(w io.WriteSeeker) *Writer {
	// Seek to the end
	offset, err := w.Seek(0, 2)
	if err != nil {
		panic("")
	}
	return &Writer{offset, nil, w}
}

// Not thread-safe
func (w *Writer) Append() (offset int64, value io.Writer, err error) {
	if w.ongoingWrite != nil {
		if err := w.ongoingWrite.finalize(); err != nil {
			return 0, nil, err
		}
	}
	offset = w.offset
	w.ongoingWrite = &recordWriter{w, w.offset, make([]byte, 0, 1000)}
	value = w.ongoingWrite
	err = nil
	return
}

func (w *Writer) Flush() error {
	if w.ongoingWrite != nil {
		if err := w.ongoingWrite.finalize(); err != nil {
			return err
		}
	}
	f, ok := w.w.(flusher)
	if ok {
		return f.Flush()
	}
	return nil
}
