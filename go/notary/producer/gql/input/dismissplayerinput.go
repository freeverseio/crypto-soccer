package input

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/helper"
)

type DismissPlayerInput struct {
	Signature       string
	PlayerId        string
	ValidUntil      string
	ReturnToAcademy bool
}

func (b DismissPlayerInput) Hash() (common.Hash, error) {
	uint256Ty, _ := abi.NewType("int256", "int256", nil)
	boolTy, _ := abi.NewType("bool", "bool", nil)

	arguments := abi.Arguments{
		{Type: uint256Ty},
		{Type: uint256Ty},
		{Type: boolTy},
	}

	validUntil, _ := new(big.Int).SetString(b.ValidUntil, 10)
	if validUntil == nil {
		return common.Hash{}, errors.New("invalid validUntil")
	}
	playerId, _ := new(big.Int).SetString(string(b.PlayerId), 10)
	if playerId == nil {
		return common.Hash{}, errors.New("invalid playerId")
	}

	bytes, err := arguments.Pack(
		validUntil,
		playerId,
		b.ReturnToAcademy,
	)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(bytes), nil
}

func (b DismissPlayerInput) SignerAddress() (common.Address, error) {
	hash, err := b.Hash()
	if err != nil {
		return common.Address{}, err
	}
	hash = helper.PrefixedHash(hash)
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return helper.AddressFromSignature(hash, sign)
}

func (b DismissPlayerInput) VerifySignature() (bool, error) {
	hash, err := b.Hash()
	if err != nil {
		return false, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return false, err
	}
	return helper.VerifySignature(hash, sign)
}

func (b DismissPlayerInput) IsSignerOwner(contracts contracts.Contracts) (bool, error) {
	signerAddress, err := b.SignerAddress()
	if err != nil {
		return false, err
	}
	playerId, _ := new(big.Int).SetString(string(b.PlayerId), 10)
	if playerId == nil {
		return false, errors.New("invalid teamId")
	}
	owner, err := contracts.Market.GetOwnerPlayer(&bind.CallOpts{}, playerId)
	if err != nil {
		return false, err
	}
	return signerAddress == owner, nil
}
