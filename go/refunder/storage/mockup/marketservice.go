package mockup

import (
	"github.com/freeverseio/crypto-soccer/go/refunder/storage"
)

type MarketService struct {
	BeginFunc func() (storage.MarketTx, error)
}

func (b *MarketService) Begin() (storage.MarketTx, error) {
	return b.BeginFunc()
}

type MarketTx struct {
	RollbackFunc func() error
	CommitFunc   func() error
}

func (b *MarketTx) Commit() error {
	return b.CommitFunc()
}
func (b *MarketTx) Rollback() error {
	return b.RollbackFunc()
}
