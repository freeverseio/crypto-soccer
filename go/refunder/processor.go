package refunder

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/refunder/storage/market"
	"github.com/freeverseio/crypto-soccer/go/refunder/storage/universe"
)

type Processor struct {
	uService universe.Service
	mService market.Service
}

func New(
	uService universe.Service,
	mService market.Service,
) (*Processor, error) {
	if uService == nil || mService == nil {
		return nil, errors.New("invalid params")
	}

	return &Processor{
		uService,
		mService,
	}, nil
}
