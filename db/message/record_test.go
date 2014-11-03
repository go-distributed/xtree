package message

import (
	"reflect"
	"testing"
)

func TestLogRecord(t *testing.T) {
	tests := []struct {
		rec *Record
	}{
		{&Record{
			Rev:  1,
			Key:  "/test",
			Data: []byte("some data"),
		}},
	}

	for i, tt := range tests {
		var data []byte
		var err error
		if data, err = tt.rec.Marshal(); err != nil {
			t.Fatalf("#%d: cannot marshal, err: %v", i, err)
		}

		rec := &Record{}
		if err = rec.Unmarshal(data); err != nil {
			t.Fatalf("#%d: cannot unmarshal, err: %v", i, err)
		}

		if !reflect.DeepEqual(tt.rec, rec) {
			t.Fatalf("#%d: records are not the same, want: %v, get: %v",
				i, tt.rec, rec)
		}
	}
}
