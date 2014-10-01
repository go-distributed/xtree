package recordio

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestFetch(t *testing.T) {
	writeFile, err := ioutil.TempFile("", "backend")
	if err != nil {
		t.Error("can't create temp file")
	}
	readFile, err := os.Open(writeFile.Name())
	if err != nil {
		t.Error("can't open temp file")
	}

	ap := NewAppender(writeFile)
	fc := NewFetcher(readFile)

	tests := []Record{
		{[]byte("someData")},
		{[]byte("someOtherData")},
	}

	for i, recWrite := range tests {
		offset, err := ap.Append(recWrite)
		if err != nil {
			t.Errorf("#%d: Append failed: %s", i, err.Error())
		}

		recRead, err := fc.Fetch(offset)
		if err != nil {
			t.Errorf("#%d: Fetch failed: %s", i, err.Error())
		}

		if !reflect.DeepEqual(recRead, recWrite) {
			t.Errorf("#%d: records not the same, want: %v, get %v",
				i, recWrite, recRead)
		}

	}
}
