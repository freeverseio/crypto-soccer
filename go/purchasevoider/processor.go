package purchasevoider

import (
	"errors"
)

type Processor struct {
	orderS    VoidPurchaseService
	universeS UniverseService
	marketS   MarketService
}

func New(
	orderS VoidPurchaseService,
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

func (b Processor) GetVoidedTokens() ([]string, error) {
	var tokens []string
	voidPurchases, err := b.orderS.VoidedPurchases(nil)
	if err != nil {
		return tokens, err
	}

	for _, p := range voidPurchases {
		tokens = append(tokens, p.PurchaseToken)
	}
	return tokens, nil
}
