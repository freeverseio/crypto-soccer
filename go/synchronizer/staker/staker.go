package staker

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go/helper"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	log "github.com/sirupsen/logrus"
)

type Staker struct {
	privateKey *ecdsa.PrivateKey
}

func New(privateKey *ecdsa.PrivateKey) (*Staker, error) {
	log.Infof("[staker] created with address : %v", crypto.PubkeyToAddress(privateKey.PublicKey).Hex())

	staker := Staker{}
	staker.privateKey = privateKey

	return &staker, nil
}

func (b Staker) Init(contracts contracts.Contracts) error {
	isTrustedParty, err := b.IsTrustedParty(contracts)
	if err != nil {
		return err
	}
	if !isTrustedParty {
		return errors.New("[staker] not a trusted party")
	}
	isEnrolled, err := b.IsEnrolled(contracts)
	if err != nil {
		return err
	}
	if !isEnrolled {
		log.Info("[staker] trying to enroll")
		stake, err := requiredStake(contracts)
		if err != nil {
			return err
		}
		log.Infof("[staker] stake required %v", stake.String())
		if err := b.enroll(contracts, stake); err != nil {
			return err
		}
		log.Info("[staker] enrollment successful")
	}
	return nil
}

func (b Staker) Address() common.Address {
	return crypto.PubkeyToAddress(b.privateKey.PublicKey)
}

func (b Staker) IsEnrolled(contracts contracts.Contracts) (bool, error) {
	return contracts.Stakers.IsStaker(&bind.CallOpts{}, b.Address())
}

func (b Staker) IsTrustedParty(contracts contracts.Contracts) (bool, error) {
	return contracts.Stakers.IsTrustedParty(&bind.CallOpts{}, b.Address())
}

func (b Staker) SubmitRoot(contracts contracts.Contracts, root [32]byte) error {
	auth := bind.NewKeyedTransactor(b.privateKey)
	auth.GasPrice = big.NewInt(1000000000) // in xdai is fixed to 1 GWei
	tx, err := contracts.Updates.UpdateTZ(auth, root)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(contracts.Client, tx, 60)
	if err != nil {
		return err
	}
	return nil
}

func (b Staker) Play(contracts contracts.Contracts, root [32]byte) error {
	isTimeToUpdate, err := contracts.Updates.IsTimeToUpdate(&bind.CallOpts{})
	if err != nil {
		return err
	}
	if !isTimeToUpdate {
		return nil
	}
	log.Infof("[staker] it's time ... let's try to submit universe root %v", hex.EncodeToString(root[:]))
	return b.SubmitRoot(contracts, root)
}

func requiredStake(contracts contracts.Contracts) (*big.Int, error) {
	stake, err := contracts.Stakers.RequiredStake(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return stake, nil
}

func (b Staker) enroll(contracts contracts.Contracts, stake *big.Int) error {
	auth := bind.NewKeyedTransactor(b.privateKey)
	auth.GasPrice = big.NewInt(1000000000) // in xdai is fixed to 1 GWei
	auth.Value = stake

	tx, err := contracts.Stakers.Enroll(auth)
	if err != nil {
		return err
	}
	if _, err := helper.WaitReceipt(contracts.Client, tx, 60); err != nil {
		return err
	}
	return nil
}
