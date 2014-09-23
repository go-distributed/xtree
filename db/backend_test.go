package db

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPut(t *testing.T) {
	tests := []struct {
		rev  int
		path string
		data []byte
	}{
		{1, "/foo/bar", []byte("somedata")},
	}

	for i, tt := range tests {
		b := newBackend()
		b.Put(tt.rev, tt.path, tt.data)
		v := b.Get(tt.rev, tt.path)
		if v.rev != tt.rev {
			t.Errorf("#%d: rev = %d, want %d", i, v.rev, tt.rev)
		}
		if !reflect.DeepEqual(v.data, tt.data) {
			t.Errorf("#%d: data = %s, want %s", i, v.data, tt.data)
		}
	}
}

func TestPutOverwrite(t *testing.T) {
	tests := []struct {
		path       string
		firstData  []byte
		secondData []byte
	}{
		{"/a", []byte("first"), []byte("second")},
	}

	for _, tt := range tests {
		b := newBackend()

		b.Put(1, tt.path, tt.firstData)
		v := b.Get(1, tt.path)
		if v.rev != 1 || !reflect.DeepEqual(v.data, tt.firstData) {
			t.Errorf("Put(1, %s) => (%d, %s)", tt.firstData, v.rev, v.data)
		}

		b.Put(2, tt.path, tt.secondData)
		v = b.Get(2, tt.path)
		if v.rev != 2 || !reflect.DeepEqual(v.data, tt.secondData) {
			t.Errorf("Put(2, %s) => (%d, %s)", tt.secondData, v.rev, v.data)
		}
	}
}

func TestLs(t *testing.T) {
	back := newBackend()
	d := []byte("somedata")
	back.Put(1, "/a", d)
	back.Put(2, "/a/b", d)
	back.Put(3, "/a/c", d)
	back.Put(4, "/b", d)

	tests := []struct {
		p   string
		wps []string
	}{
		{"/", []string{"/a", "/b"}},
		{"/a", []string{"/a/b", "/a/c"}},
		{"/a/", []string{"/a/b", "/a/c"}},
		{"/a/b", []string{}},
		{"/b", []string{}},
		{"/c", []string{}},
	}
	for i, tt := range tests {
		ps := back.Ls(tt.p)
		if len(ps) != len(tt.wps) {
			t.Fatalf("#%d: len(ps) = %d, want %d", i, len(ps), len(tt.wps))
		}
		for j := range ps {
			if ps[j].p != tt.wps[j] {
				t.Errorf("#%d.%d: path = %s, want %s", i, j, ps[j].p, tt.wps[j])
			}
		}
	}
}

func BenchmarkPut(b *testing.B) {
	b.StopTimer()
	back := newBackend()
	d := []byte("somedata")
	path := make([]string, b.N)
	for i := range path {
		path[i] = fmt.Sprintf("/%d", i+1)
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
	path := make([]string, b.N)
	for i := range path {
		path[i] = fmt.Sprintf("/%d", i+1)
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
	path := make([]string, b.N)
	for i := range path {
		path[i] = fmt.Sprintf("/%d", i+1)
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
