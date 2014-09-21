package db

// Batch is a collection of operations that's expected to be executed atomically.
type Batch struct {
	batchOps []*batchOp
}

type batchOp struct {
}

// add a put operation to this batch
func (b *Batch) Put(rev int, path string, data []byte) error {

}
