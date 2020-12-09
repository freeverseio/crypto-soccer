package refunder

import (
	"errors"
)

type Processor struct {
	pService PaymentService
	uService UniverseService
	mService MarketService
}

func New(
	pService PaymentService,
	uService UniverseService,
	mService MarketService,
) (*Processor, error) {
	if uService == nil || mService == nil {
		return nil, errors.New("invalid params")
	}

	return &Processor{
		pService,
		uService,
		mService,
	}, nil
}

func GetVoidedOrders() error {
	return nil
}
