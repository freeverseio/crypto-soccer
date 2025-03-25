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

func DismissPlayer(
	contracts contracts.Contracts,
	pvc *ecdsa.PrivateKey,
	in input.DismissPlayerInput,
) error {
	validUntil, _ := new(big.Int).SetString(in.ValidUntil, 10)
	if validUntil == nil {
		return errors.New("invalid validUntil")
	}
	playerId, _ := new(big.Int).SetString(in.PlayerId, 10)
	if playerId == nil {
		return errors.New("invalid playerId")
	}
	r, s, v, err := helper.RSV(in.Signature)
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(pvc)
	auth.GasPrice = big.NewInt(3000000000) // in xdai is fixe to 3 GWei
	tx, err := contracts.Market.DismissPlayer(
		auth,
		validUntil,
		playerId,
		r,
		s,
		v,
	)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceiptAndCheckSuccess(contracts.Client, tx, 30)
	if err != nil {
		return err
	}

	return nil
}
