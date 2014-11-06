package db

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	dblog "github.com/go-distributed/xtree/db/log"
	"github.com/go-distributed/xtree/db/message"
	"github.com/go-distributed/xtree/third-party/github.com/google/btree"
)

type backend struct {
	bt     *btree.BTree
	cache  *cache
	rev    int
	dblog  *dblog.DBLog
	config *DBConfig
}

func newBackend() *backend {
	dataDir, err := ioutil.TempDir("", "backend")
	if err != nil {
		panic("not implemented")
	}

	config := &DBConfig{
		DataDir: dataDir,
	}
	b, err := newBackendWithConfig(config)
	if err != nil {
		panic("not implemented")
	}
	return b
}

func newBackendWithConfig(config *DBConfig) (b *backend, err error) {
	bt := btree.New(10)
	b = &backend{
		bt:     bt,
		cache:  newCache(),
		config: config,
	}
	haveLog := dblog.Exist(config.DataDir)
	switch haveLog {
	case false:
		fmt.Println("didn't have log file. Init...")
		err = b.init(config)
	case true:
		fmt.Println("had log file. Restore...")
		err = b.restore(config)
	}
	return
}

func (b *backend) getData(offset int64) []byte {
	rec, err := b.dblog.GetRecord(offset)
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
	offset, err := b.dblog.Append(&message.Record{
		Key:  path.p,
		Data: data,
	})
	if err != nil {
		panic("unimplemented")
	}

	nv.offset = offset
}

// one-level listing
func (b *backend) Ls(pathname string) (paths []Path) {
	paths = make([]Path, 0)
	pivot := newPathForLs(pathname)

	b.bt.AscendGreaterOrEqual(pivot, func(treeItem btree.Item) bool {
		p := treeItem.(*Path)
		if !strings.HasPrefix(p.p, pivot.p) ||
			p.level != pivot.level {
			return false
		}
		paths = append(paths, *p)
		return true
	})

	return
}

// init() creates a new log file
func (b *backend) init(config *DBConfig) (err error) {
	b.dblog, err = dblog.Create(config.DataDir)
	return
}

// restore() restores database from the log file.
func (b *backend) restore(config *DBConfig) (err error) {
	rev := 0
	return dblog.Reuse(config.DataDir,
		func(l *dblog.DBLog) {
			b.dblog = l
		},
		func(r *message.Record) (err error) {
			rev++
			p := newPath(r.Key)
			b.Put(rev, *p, r.Data)
			return
		})
}

// clean up resource after testing
func (b *backend) testableCleanupResource() (err error) {
	b.dblog.Close()
	return os.RemoveAll(b.config.DataDir)
}
