package log

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/go-distributed/xtree/db/message"
)

func TestAppendAndGetRecord(t *testing.T) {
	var err error
	var l *DBLog
	var dataDir string

	if dataDir, err = ioutil.TempDir("", "logtest"); err != nil {
		t.Errorf("ioutil.TempDir failed: %v", err)
	}

	defer os.RemoveAll(dataDir)

	if l, err = Create(dataDir); err != nil {
		t.Errorf("Create failed: %v", err)
	}

	tests := []struct {
		offset int64
		record *message.Record
	}{
		{-1, &message.Record{
			Key:  "/test",
			Data: []byte("some data"),
		}},
		{-1, &message.Record{
			Key:  "/test2",
			Data: []byte("some other data"),
		}},
	}

	for i, tt := range tests {
		tests[i].offset, err = l.Append(tt.record)
		if err != nil {
			t.Errorf("#%d: Append failed: %v", i, err)
		}
	}

	for i, tt := range tests {
		var r *message.Record
		if r, err = l.GetRecord(tt.offset); err != nil {
			t.Errorf("#%d: GetRecord failed: %v", i, err)
		}

		if !reflect.DeepEqual(tt.record, r) {
			t.Errorf("#%d: records not the same, want: %v, get %v",
				i, tt.record, r)
		}

	}
}
