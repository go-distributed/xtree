package db

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-distributed/xtree/db/recordio"
	"github.com/go-distributed/xtree/third-party/github.com/google/btree"
)

type backend struct {
	bt    *btree.BTree
	cache *cache
	rev   int
	fc    recordio.Fetcher
	ap    recordio.Appender
}

func newBackend() *backend {
	bt := btree.New(10)

	// temporary file IO to test in-disk values
	writeFile, err := ioutil.TempFile("", "backend")
	if err != nil {
		panic("can't create temp file")
	}
	readFile, err := os.Open(writeFile.Name())
	if err != nil {
		panic("can't open temp file")
	}

	return &backend{
		bt:    bt,
		cache: newCache(),
		fc:    recordio.NewFetcher(readFile),
		ap:    recordio.NewAppender(writeFile),
	}
}

func (b *backend) getData(offset int64) []byte {
	rec := &recordio.Record{Data: nil}
	err := b.fc.Fetch(offset, rec)
	if err != nil {
		panic("unimplemented")
	}
	return rec.Data
}

// if it couldn't find anything related to path, it return Value of 0 rev.
func (b *backend) Get(rev int, path Path) Value {
	if b.cache != nil {
		if v, ok := b.cache.get(revpath{rev: rev, path: path.p}); ok {
			return v
		}
	}
	item := b.bt.Get(&path)
	if item == nil {
		return Value{}
	}

	v := item.(*Path).v

	for v != nil && v.rev > rev {
		v = v.next
	}

	if v == nil {
		return Value{}
	}

	return Value{
		rev:  v.rev,
		data: b.getData(v.offset),
	}
}

func (b *backend) Put(rev int, path Path, data []byte) {
	nv := &memValue{rev: b.rev + 1}
	item := b.bt.Get(&path)
	if item == nil {
		path.v = nv
		b.bt.ReplaceOrInsert(&path)
	} else {
		exPath := item.(*Path)
		nv.next = exPath.v
		exPath.v = nv
	}

	b.rev++
	offset, err := b.ap.Append(&recordio.Record{data})
	if err != nil {
		panic("unimplemented")
	}

	nv.offset = offset
}

// one-level listing
func (b *backend) Ls(pathname string) []Path {
	result := make([]Path, 0)
	pivot := newPathForLs(pathname)

	b.bt.AscendGreaterOrEqual(pivot, func(treeItem btree.Item) bool {
		p := treeItem.(*Path)
		if !strings.HasPrefix(p.p, pivot.p) ||
			p.level != pivot.level {
			return false
		}
		result = append(result, *p)
		return true
	})
	return result
}
