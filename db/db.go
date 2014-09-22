package db

import "github.com/go-distributed/xtree/third-party/github.com/google/btree"

// DB is the data store.
type DB interface {
	// Get returns the value of path, only if path exists and rev is no smaller than
	// the last rev of path or rev is -1; Otherwise an error is returned.
	Get(rev int, path string) (Value, error)
	// Put sets data on path (if a nil data is set, the specific path is
	// deleted), only if path exists and rev is no smaller than
	// the last rev of path or rev is -1.
	// If data is set successfully, the value of path is returned.
	// Otherwise, an error is returned.
	Put(rev int, path string, data []byte) (Value, error)
	// Head returns the global rev of this DB.
	// If it failed, an error is returned.
	Head() (int, error)
	// List returns a possibly-empty list of Paths with prefix at or after rev.
	// The list is sorted in the ordering of:
	// 1. level (separated by '/')
	// 2. name in each level.
	//
	// If level is >= 0, it takes up to level from current path. '/a' is
	// 0 level from '/a', while '/a/b' is 1.
	// If level is -1, it means any.
	//
	// if count is >= 0, it is how many paths we want in the list.
	// if count is -1, it means any.
	//
	// if it failed, an error is returned.
	List(rev int, prefix string, level int, count int) ([]Path, error)
	// Transaction executes a batch of operations atomically.
	// If it succeeds, a list of values corresponding to each operation is
	// returned. Otherwise, an error indicating the first failure is returned.
	Transaction(batch *Batch) ([]Value, error)
}

// Value is the value on specific path.
type Value struct {
	rev  int
	next *Value
	data []byte // todo: file offset
}

// Path is a collection of infomation on specific path.
type Path struct {
	p string
	v Value
}

func (a Path) Less(b btree.Item) bool {
	return a.p < b.(Path).p
}
