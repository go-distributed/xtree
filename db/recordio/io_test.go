package recordio

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestFetch(t *testing.T) {
	writeFile, err := ioutil.TempFile("", "testfetch")
	if err != nil {
		t.Error("can't create temp file")
	}
	defer os.Remove(writeFile.Name())
	defer writeFile.Close()

	readFile, err := os.Open(writeFile.Name())
	if err != nil {
		t.Error("can't open temp file")
	}
	defer readFile.Close()

	ap := NewAppender(writeFile)
	fc := NewFetcher(readFile)

	tests := []struct {
		offset int64
		record Record
	}{
		{-1, Record{[]byte("someData")}},
		{-1, Record{[]byte("someOtherData")}},
	}

	for i, tt := range tests {
		offset, err := ap.Append(tt.record)
		if err != nil {
			t.Errorf("#%d: Append failed: %s", i, err.Error())
		}
		tests[i].offset = offset
	}

	for i, tt := range tests {
		recRead, err := fc.Fetch(tt.offset)
		if err != nil {
			t.Errorf("#%d: Fetch failed: %s", i, err.Error())
		}

		if !reflect.DeepEqual(recRead, tt.record) {
			t.Errorf("#%d: records not the same, want: %v, get %v",
				i, tt.offset, recRead)
		}

	}
}
