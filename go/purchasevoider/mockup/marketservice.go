package mockup

import (
	"github.com/freeverseio/crypto-soccer/go/purchasevoider"
)

type MarketService struct {
	BeginFunc func() (purchasevoider.MarketTx, error)
}

func (b *MarketService) Begin() (purchasevoider.MarketTx, error) {
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
