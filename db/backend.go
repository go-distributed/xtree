package db

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-distributed/xtree/db/record"
	"github.com/go-distributed/xtree/third-party/github.com/google/btree"
)

type backend struct {
	bt     *btree.BTree
	cache  *cache
	rev    int
	reader *record.Reader
	writer *record.Writer
}

func newBackend() *backend {
	bt := btree.New(10)

	file, err := os.OpenFile("test-records",
		os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_TRUNC,
		os.FileMode(0644))
	if err != nil {
		panic("can't open file")
	}

	return &backend{
		bt:     bt,
		cache:  newCache(),
		reader: record.NewReader(file),
		writer: record.NewWriter(file),
	}
}

func (b *backend) Close() {
	os.Remove("test-records")
}

func (b *backend) getData(offset int64) []byte {
	r, err := b.reader.ReadAt(offset)
	if err != nil {
		panic("unimplemented")
	}

	bs, err := ioutil.ReadAll(r)
	if err != nil {
		panic("unimplemented")
	}

	return bs
}

// if it couldn't find anything related to path, it return Value of 0 rev.
func (b *backend) Get(rev int, path Path) Value {
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

	// offset, data
	offset, w, err := b.writer.Append()
	if err != nil {
		panic("unimplemented")
	}

	nv.offset = offset
	_, err = w.Write(data)
	if err != nil {
		panic("unimplemented")
	}

	b.writer.Flush()
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
