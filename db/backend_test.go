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
		{1, *newPath("/foo/bar"), []byte("somedata")},
		{2, *newPath("/bar/foo"), []byte("datasome")},
	}

	b := newBackend()
	for i, tt := range tests {
		b.Put(tt.rev, tt.path, tt.data)
		v := b.Get(tt.rev, tt.path)
		if v.rev != tt.rev {
			t.Errorf("#%d: rev = %d, want %d", i, v.rev, tt.rev)
		}
		if !reflect.DeepEqual(v.data, tt.data) {
			t.Errorf("#%d: data = %s, want %s", i, v.data, tt.data)
		}
	}
	b.Close()
}

func TestPutOnExistingPath(t *testing.T) {
	tests := []struct {
		path  Path
		data1 []byte
		data2 []byte
	}{
		{*newPath("/foo/bar"), []byte("first"), []byte("second")},
		{*newPath("/bar/foo"), []byte("first"), []byte("second")},
	}

	b := newBackend()
	for i, tt := range tests {
		b.Put(2*i+1, tt.path, tt.data1)
		v := b.Get(2*i+1, tt.path)

		if v.rev != 2*i+1 {
			t.Errorf("#%d 1: rev = %d, want %d", i, v.rev, 2*i+1)
		}
		if !reflect.DeepEqual(v.data, tt.data1) {
			t.Errorf("#%d 1: data = %s, want %s", i, v.data, tt.data1)
		}

		b.Put(2*i+2, tt.path, tt.data2)
		v = b.Get(2*i+2, tt.path)

		if v.rev != 2*i+2 {
			t.Errorf("#%d 2: rev = %d, want %d", i, v.rev, 2*i+2)
		}
		if !reflect.DeepEqual(v.data, tt.data2) {
			t.Errorf("#%d 2: data = %s, want %s", i, v.data, tt.data2)
		}
	}
	b.Close()
}

func TestGetMVCC(t *testing.T) {
	b := newBackend()
	b.Put(1, *newPath("/a"), []byte("1"))
	b.Put(2, *newPath("/b"), []byte("2"))
	b.Put(3, *newPath("/a"), []byte("3"))
	b.Put(4, *newPath("/b"), []byte("4"))

	tests := []struct {
		getRev  int
		wantRev int
		path    Path
		data    []byte
	}{
		{1, 1, *newPath("/a"), []byte("1")},
		{2, 1, *newPath("/a"), []byte("1")},
		{3, 3, *newPath("/a"), []byte("3")},
		{1, 0, *newPath("/b"), nil},
		{2, 2, *newPath("/b"), []byte("2")},
		{3, 2, *newPath("/b"), []byte("2")},
		{4, 4, *newPath("/b"), []byte("4")},
	}

	for i, tt := range tests {
		v := b.Get(tt.getRev, tt.path)

		if v.rev != tt.wantRev {
			t.Errorf("#%d: rev = %d, want %d", i, v.rev, tt.wantRev)
		}

		if !reflect.DeepEqual(v.data, tt.data) {
			t.Errorf("#%d: data = %s, want %s", i, v.data, tt.data)
		}
	}
	b.Close()
}

func TestLs(t *testing.T) {
	back := newBackend()
	d := []byte("somedata")
	back.Put(1, *newPath("/a"), d)
	back.Put(2, *newPath("/a/b"), d)
	back.Put(3, *newPath("/a/c"), d)
	back.Put(4, *newPath("/b"), d)

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
	back.Close()
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
	back.Close()
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
	back.Close()
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
	back.Close()
}
