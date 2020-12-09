package mockup

import (
	"github.com/freeverseio/crypto-soccer/go/refunder/storage"
)

type UniverseService struct {
	BeginFunc func() (storage.UniverseTx, error)
}

func (b *UniverseService) Begin() (storage.UniverseTx, error) {
	return b.BeginFunc()
}

type UniverseTx struct {
	RollbackFunc func() error
	CommitFunc   func() error
}

func (b *UniverseTx) Commit() error {
	return b.CommitFunc()
}
func (b *UniverseTx) Rollback() error {
	return b.RollbackFunc()
}
