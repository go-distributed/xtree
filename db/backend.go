package db

import (
	"fmt"
	"strings"

	"github.com/go-distributed/xtree/third-party/github.com/google/btree"
)

type backend struct {
	bt    *btree.BTree
	cache *cache
	rev   int
}

func newBackend() *backend {
	bt := btree.New(10)
	return &backend{
		bt:    bt,
		cache: newCache(),
	}
}

func (b *backend) Get(rev int, path Path) Value {
	if b.cache != nil {
		if v, ok := b.cache.get(revpath{rev: rev, path: path.p}); ok {
			return v
		}
	}
	item := b.bt.Get(&path)
	if item == nil {
		panic("unimplemented")
	}

	p := item.(*Path)
	if p.v.rev == rev {
		return *p.v
	}
	fmt.Println(p.v.rev, rev)
	panic("unimplemented")
}

func (b *backend) Put(rev int, path Path, data []byte) {
	nv := &Value{
		rev:  b.rev + 1,
		data: data,
	}
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
