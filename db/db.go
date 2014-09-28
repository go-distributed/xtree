package db

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
	// Ls returns a possibly-empty list of Paths under speicifc path at or after rev.
	// The list is sorted by the ordering of:
	// 1. level (separated by '/')
	// 2. name in each level.
	//
	// If recursive is false, it lists just one level under the path.
	// Otherwise, it lists recursively all paths.
	//
	// if count is >= 0, it is the number of paths we want in the list.
	// if count is -1, it means any.
	//
	// if it failed, an error is returned.
	Ls(rev int, path string, recursive bool, count int) ([]Path, error)
	// Transaction executes a batch of operations atomically.
	// If it succeeds, a list of values corresponding to each operation is
	// returned. Otherwise, an error indicating the first failure is returned.
	Transaction(batch *Batch) ([]Value, error)
}

// Value is the value on specific path.
type Value struct {
	rev  int
	data []byte // todo: file offset
}

type memValue struct {
	rev    int
	next   *memValue
	offset int64
}
