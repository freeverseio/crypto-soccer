package mockup

import "github.com/freeverseio/crypto-soccer/go/refunder/storage/universe"

type Service struct {
	BeginFunc func() (universe.Tx, error)
}

func (b *Service) Begin() (universe.Tx, error) {
	return b.BeginFunc()
}
