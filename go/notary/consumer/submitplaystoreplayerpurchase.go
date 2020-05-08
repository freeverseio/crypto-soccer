package consumer

import (
	"fmt"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
)

func SubmitPlayStorePlayerPurchase(in input.SubmitPlayStorePlayerPurchaseInput) error {
	playerId, _ := new(big.Int).SetString(string(in.PlayerId), 10)
	if playerId == nil {
		return fmt.Errorf("invalid playerId %v", in.PlayerId)
	}
	teamId, _ := new(big.Int).SetString(string(in.TeamId), 10)
	if teamId == nil {
		return fmt.Errorf("invalid teamId %v", in.TeamId)
	}

	return nil
}
