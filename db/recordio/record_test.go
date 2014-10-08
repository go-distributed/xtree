package recordio

import (
	"bytes"
	"reflect"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		data []byte
	}{
		{[]byte("someData")},
		{[]byte("someOtherData")},
	}

	for i, tt := range tests {
		buf := new(bytes.Buffer)
		recordToWrite := &Record{tt.data}

		if err := recordToWrite.encodeTo(buf); err != nil {
			t.Fatalf("#%d: cannot encode, err: %s", i, err)
		}
		recordToRead := new(Record)
		if err := recordToRead.decodeFrom(buf); err != nil {
			t.Fatalf("#%d: cannot decode, err: %s", i, err)
		}
		if !reflect.DeepEqual(recordToRead, recordToWrite) {
			t.Fatalf("#%d: records are not the same, want: %v, get: %v",
				i, recordToWrite, recordToRead)
		}
	}
}
