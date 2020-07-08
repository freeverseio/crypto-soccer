package consumer

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
)

func CompletePlayerTransit(
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	in input.CompletePlayerTransitInput,
) error {
	playerId, _ := new(big.Int).SetString(string(in.PlayerId), 10)
	if playerId == nil {
		return errors.New("invalid playerId")
	}

	auth := bind.NewKeyedTransactor(pvc)
	auth.GasPrice = big.NewInt(1000000000) // in xdai is fixe to 1 GWei
	tx, err := contracts.Market.CompletePlayerTransit(
		auth,
		playerId,
	)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(contracts.Client, tx, 30)
	if err != nil {
		return err
	}

	return nil
}
