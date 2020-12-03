package consumer

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql"
)

type ConsumePromo struct {
	client    *ethclient.Client
	contracts contracts.Contracts
	auth      *bind.TransactOpts
}

func NewConsumePromo(
	client *ethclient.Client,
	auth *bind.TransactOpts,
	contracts contracts.Contracts,
) *ConsumePromo {
	return &ConsumePromo{
		client,
		contracts,
		auth,
	}
}

func (b ConsumePromo) Process(in gql.ConsumePromoInput) error {
	return b.assignAsset(in)
}

func (b ConsumePromo) assignAsset(in gql.ConsumePromoInput) error {
	playerId, _ := new(big.Int).SetString(in.PlayerId, 10)
	if playerId == nil {
		return errors.New("invalid player")
	}
	teamId, _ := new(big.Int).SetString(in.TeamId, 10)
	if teamId == nil {
		return errors.New("invalid team")
	}

	tx, err := b.contracts.Market.TransferBuyNowPlayer(
		b.auth,
		playerId,
		teamId,
	)
	if err != nil {
		return err
	}
	if _, err = helper.WaitReceipt(b.contracts.Client, tx, 60); err != nil {
		return err
	}
	return nil
}
