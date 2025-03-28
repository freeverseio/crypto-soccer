package consumer

import (
	"fmt"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type FirstBotTransfer struct {
	client    *ethclient.Client
	contracts contracts.Contracts
	auth      *bind.TransactOpts
}

func NewFirstBotTransfer(
	client *ethclient.Client,
	auth *bind.TransactOpts,
	contracts contracts.Contracts,
) *FirstBotTransfer {
	return &FirstBotTransfer{
		client,
		contracts,
		auth,
	}
}

func (b FirstBotTransfer) Process(event gql.TransferFirstBotToAddrInput) error {
	timezone := uint8(event.Timezone)
	countryIdxInTimezone, _ := new(big.Int).SetString(event.CountryIdxInTimezone, 10)
	if countryIdxInTimezone == nil {
		return fmt.Errorf("Invalid countryIdxInTimezone %v", event.CountryIdxInTimezone)
	}
	address := common.HexToAddress(event.Address)
	transaction, err := b.contracts.Assets.TransferFirstBotToAddr(b.auth, timezone, countryIdxInTimezone, address)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceiptAndCheckSuccess(b.client, transaction, 30) // TODO make timeout configurable
	if err != nil {
		return err
	}
	return nil
}
