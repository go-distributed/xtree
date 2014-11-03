package log

import (
	"reflect"
	"testing"

	"github.com/go-distributed/xtree/db/message"
)

func TestAppendAndGetRecord(t *testing.T) {
	var err error
	var log *Log
	if log, err = Create(); err != nil {
		t.Errorf("Create failed: %v", err)
	}

	defer log.Destroy()

	tests := []struct {
		offset int64
		rec    *message.Record
	}{
		{-1, &message.Record{
			Rev:  1,
			Key:  "/test",
			Data: []byte("some data"),
		}},
		{-1, &message.Record{
			Rev:  2,
			Key:  "/test2",
			Data: []byte("some other data"),
		}},
	}

	for i, tt := range tests {
		tests[i].offset, err = log.Append(tt.rec)
		if err != nil {
			t.Errorf("#%d: Append failed: %v", i, err)
		}
	}

	for i, tt := range tests {
		var rec *message.Record
		if rec, err = log.GetRecord(tt.offset); err != nil {
			t.Errorf("#%d: GetRecord failed: %v", i, err)
		}

		if !reflect.DeepEqual(tt.rec, rec) {
			t.Errorf("#%d: records not the same, want: %v, get %v",
				i, tt.rec, rec)
		}

	}
}
