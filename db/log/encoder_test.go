package log

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/go-distributed/xtree/db/message"
)

func TestEncoderDecoder(t *testing.T) {
	tests := []struct {
		rec *message.Record
	}{
		{&message.Record{
			Key:  "/test",
			Data: []byte("some data"),
		}},
	}

	for i, tt := range tests {
		var err error
		eBuf := new(bytes.Buffer)
		encoder := newEncoder(eBuf)

		if err = encoder.encode(tt.rec); err != nil {
			t.Fatalf("#%d: cannot encode, err: %v", i, err)
		}
		if err = encoder.flush(); err != nil {
			t.Fatalf("#%d: cannot flush encode, err: %v", i, err)
		}

		rec := &message.Record{}
		dBuf := bytes.NewBuffer(eBuf.Bytes())
		decoder := newDecoder(dBuf)

		if err = decoder.decode(rec); err != nil {
			t.Fatalf("#%d: cannot decode, err: %v", i, err)
		}

		if !reflect.DeepEqual(tt.rec, rec) {
			t.Fatalf("#%d: records are not the same, want: %v, get: %v",
				i, tt.rec, rec)
		}
	}
}
