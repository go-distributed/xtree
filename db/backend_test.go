package db

import (
	"reflect"
	"testing"
)

func TestPut(t *testing.T) {
	tests := []struct {
		rev  int
		path Path
		data []byte
	}{
		{1, Path{p: "/foo/bar"}, []byte("somedata")},
	}

	for i, tt := range tests {
		b := newBackend()
		b.Put(tt.rev, tt.path, tt.data)
		v := b.Get(tt.rev, tt.path)
		if v.rev != tt.rev {
			t.Errorf("#%d: rev = %d, want %d", i, v.rev, tt.rev)
		}
		if !reflect.DeepEqual(v.data, tt.data) {
			t.Errorf("#%d: data = %d, want %d", i, v.data, tt.data)
		}
	}
}
