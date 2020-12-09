package mockup

import (
	"github.com/freeverseio/crypto-soccer/go/refunder"
)

type MarketService struct {
	BeginFunc func() (refunder.MarketTx, error)
}

func (b *MarketService) Begin() (refunder.MarketTx, error) {
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
