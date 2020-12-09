package mockup

import "github.com/freeverseio/crypto-soccer/go/refunder/storage/market"

type Service struct {
	BeginFunc func() (market.Tx, error)
}

func (b *Service) Begin() (market.Tx, error) {
	return b.BeginFunc()
}
