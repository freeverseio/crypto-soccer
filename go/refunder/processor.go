package refunder

import (
	"errors"

	"github.com/freeverseio/crypto-soccer/go/refunder/storage"
)

type Processor struct {
	uService storage.UniverseService
	mService storage.MarketService
}

func New(
	uService storage.UniverseService,
	mService storage.MarketService,
) (*Processor, error) {
	if uService == nil || mService == nil {
		return nil, errors.New("invalid params")
	}

	return &Processor{
		uService,
		mService,
	}, nil
}
