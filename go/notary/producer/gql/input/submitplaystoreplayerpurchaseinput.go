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
	"github.com/graph-gophers/graphql-go"
)

type SubmitPlayStorePlayerPurchaseInput struct {
	Signature string
	Receipt   string
	PlayerId  graphql.ID
	TeamId    graphql.ID
}

func (b SubmitPlayStorePlayerPurchaseInput) Hash() (common.Hash, error) {
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	stringTy, _ := abi.NewType("string", "string", nil)

	arguments := abi.Arguments{
		{Type: stringTy},
		{Type: uint256Ty},
		{Type: uint256Ty},
	}

	teamId, _ := new(big.Int).SetString(string(b.TeamId), 10)
	if teamId == nil {
		return common.Hash{}, errors.New("Invalid TeamId")
	}
	playerId, _ := new(big.Int).SetString(string(b.PlayerId), 10)
	if playerId == nil {
		return common.Hash{}, errors.New("Invalid PlayerId")
	}

	bytes, err := arguments.Pack(
		b.Receipt,
		playerId,
		teamId,
	)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(bytes), nil
}

func (b SubmitPlayStorePlayerPurchaseInput) SignerAddress() (common.Address, error) {
	hash, err := b.Hash()
	if err != nil {
		return common.Address{}, err
	}
	hash = helper.PrefixedHash(hash)
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return helper.AddressFromHashAndSignature(hash, sign)
}

func (b SubmitPlayStorePlayerPurchaseInput) IsSignerOwner(contracts contracts.Contracts) (bool, error) {
	signerAddress, err := b.SignerAddress()
	if err != nil {
		return false, err
	}
	teamId, _ := new(big.Int).SetString(string(b.TeamId), 10)
	if teamId == nil {
		return false, errors.New("Invalid teamId")
	}
	owner, err := contracts.Market.GetOwnerTeam(&bind.CallOpts{}, teamId)
	if err != nil {
		return false, err
	}
	return signerAddress == owner, nil
}
