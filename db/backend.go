package db

import (
	"fmt"

	"github.com/google/btree"
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
	item := b.bt.Get(path)
	p := item.(Path)
	if p.v.rev == rev {
		return p.v
	}
	fmt.Println(p.v.rev, rev)
	panic("unimplemented")
}

func (b *backend) Put(rev int, path Path, data []byte) {
	nv := Value{
		rev:  b.rev + 1,
		data: data,
	}
	item := b.bt.Get(path)
	if item == nil {
		path.v = nv
		b.bt.ReplaceOrInsert(path)
		b.rev++
		return
	}
	panic("unimplemented")
}
