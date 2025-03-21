package purchasevoider

import (
	"errors"

	log "github.com/sirupsen/logrus"
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

func (b Processor) Run() error {
	log.Info("get voided tokens ...")
	tokens, err := b.GetVoidedTokens()
	if err != nil {
		return err
	}
	log.Infof("found %v tokens", len(tokens))
	playerIds, err := b.GetPlayerIds(tokens)
	if err != nil {
		return err
	}
	return b.MarkForDeletion(playerIds)
}

func (b Processor) GetVoidedTokens() ([]string, error) {
	var tokens []string
	voidPurchases, err := b.orderS.VoidedPurchases(nil) // TODO context
	if err != nil {
		return tokens, err
	}

	for _, p := range voidPurchases {
		tokens = append(tokens, p.PurchaseToken)
	}
	return tokens, nil
}

func (b Processor) GetPlayerIds(tokens []string) ([]string, error) {
	var ids []string
	for _, t := range tokens {
		id, err := b.marketS.GetPlayerIdByPurchaseToken(t)
		if err != nil {
			return ids, err
		}
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (b Processor) MarkForDeletion(playerIds []string) error {
	for _, id := range playerIds {
		if err := b.universeS.MarkForDeletion(id); err != nil {
			return err
		}
	}
	return nil
}
