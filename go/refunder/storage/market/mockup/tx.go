package mockup

type Tx struct {
	RollbackFunc func() error
	CommitFunc   func() error
}

func (b *Tx) Commit() error {
	return b.CommitFunc()
}
func (b *Tx) Rollback() error {
	return b.RollbackFunc()
}
