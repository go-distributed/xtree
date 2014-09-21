package db

import "github.com/google/btree"

type backend struct {
	bt  *btree.BTree
	rev int
}

func newBackend() *backend {
	bt := btree.New(10)
	return &backend{
		bt: bt,
	}
}

func (b *backend) Get(rev int, path Path) Value {
	item := b.bt.Get(path)
	p := item.(Path)
	if p.v.rev == rev {
		return p.v
	}
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
		return
	}
	panic("unimplemented")
}
