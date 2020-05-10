package staker

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go/helper"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	log "github.com/sirupsen/logrus"
)

type Staker struct {
	auth *bind.TransactOpts
}

func New(privateKey *ecdsa.PrivateKey) (*Staker, error) {
	log.Infof("[staker] created with address : %v", crypto.PubkeyToAddress(privateKey.PublicKey).Hex())

	staker := Staker{}
	staker.auth = bind.NewKeyedTransactor(privateKey)
	staker.auth.GasPrice = big.NewInt(1000000000) // in xdai is fixed to 1 GWei

	return &staker, nil
}

func (b Staker) Address() common.Address {
	return b.auth.From
}

func (b Staker) IsEnrolled(contracts contracts.Contracts) (bool, error) {
	return contracts.Stakers.IsStaker(&bind.CallOpts{}, b.Address())
}

func (b Staker) Enroll(contracts contracts.Contracts) error {
	tx, err := contracts.Stakers.Enroll(b.auth)
	if err != nil {
		return err
	}
	if _, err := helper.WaitReceipt(contracts.Client, tx, 60); err != nil {
		return err
	}
	return nil
}
