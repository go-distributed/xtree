package db

import (
	"fmt"
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

func BenchmarkPut(b *testing.B) {
	b.StopTimer()
	back := newBackend()
	d := []byte("somedata")
	path := make([]Path, b.N)
	for i := range path {
		path[i] = Path{p: fmt.Sprintf("/%d", i+1)}
	}

	b.StartTimer()
	for i := 1; i < b.N; i++ {
		back.Put(i, path[i], d)
	}
}

func BenchmarkGetWithCache(b *testing.B) {
	b.StopTimer()
	back := newBackend()
	d := []byte("somedata")
	path := make([]Path, b.N)
	for i := range path {
		path[i] = Path{p: fmt.Sprintf("/%d", i+1)}
	}
	for i := 1; i < b.N; i++ {
		back.Put(i, path[i], d)
	}

	b.StartTimer()
	for i := 1; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			back.Get(i, path[i])
		}
	}
}

func BenchmarkGetWithOutCache(b *testing.B) {
	b.StopTimer()
	back := newBackend()
	back.cache = nil
	d := []byte("somedata")
	path := make([]Path, b.N)
	for i := range path {
		path[i] = Path{p: fmt.Sprintf("/%d", i+1)}
	}
	for i := 1; i < b.N; i++ {
		back.Put(i, path[i], d)
	}

	b.StartTimer()
	for i := 1; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			back.Get(i, path[i])
		}
	}
}
