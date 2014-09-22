package db

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-distributed/xtree/third-party/github.com/google/btree"
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

func ExampleNewPath() {
	p := newPath("/a/b/c/")
	fmt.Println(p.level, p.p)

	p = newPathForLs("/a/b/c")
	fmt.Println(p.level, p.p)

	// Output:
	// 3 /a/b/c
	// 4 /a/b/c/
}

func ExampleSortOrder() {
	back := newBackend()
	d := []byte("somedata")
	back.Put(1, newPath("/a"), d)
	back.Put(2, newPath("/a/b"), d)
	back.Put(3, newPath("/a/c"), d)
	back.Put(4, newPath("/b"), d)

	back.bt.Ascend(func(i btree.Item) bool {
		fmt.Println(i.(Path).p)
		return true
	})
	// Output:
	// /a
	// /b
	// /a/b
	// /a/c
}

func ExampleLs() {
	back := newBackend()
	d := []byte("somedata")
	back.Put(1, newPath("/a"), d)
	back.Put(2, newPath("/a/b"), d)
	back.Put(3, newPath("/b"), d)

	ps := back.Ls("/")
	for _, p := range ps {
		fmt.Print(p.p, " ")
	}

	// Output:
	// /a /b
}

func TestLs(t *testing.T) {
	back := newBackend()
	d := []byte("somedata")
	back.Put(1, newPath("/a"), d)
	back.Put(2, newPath("/a/b"), d)
	back.Put(3, newPath("/a/c"), d)
	back.Put(4, newPath("/b"), d)

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
