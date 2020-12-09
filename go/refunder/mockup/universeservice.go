package mockup

import (
	"github.com/freeverseio/crypto-soccer/go/refunder"
)

type UniverseService struct {
	BeginFunc func() (refunder.UniverseTx, error)
}

func (b *UniverseService) Begin() (refunder.UniverseTx, error) {
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
