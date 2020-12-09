package refunder

import (
	"errors"
)

type Processor struct {
	orderS    OrderService
	universeS UniverseService
	marketS   MarketService
}

func New(
	orderS OrderService,
	universeS UniverseService,
	marketS MarketService,
) (*Processor, error) {
	if universeS == nil || marketS == nil || orderS == nil {
		return nil, errors.New("invalid params")
	}

	return &Processor{
		orderS,
		universeS,
		marketS,
	}, nil
}

func GetVoidedOrders() error {
	return nil
}
