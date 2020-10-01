package input

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/graph-gophers/graphql-go"
)

type CancelAllOffersBySellerInput struct {
	Signature string
	PlayerId  graphql.ID
}

func (b CancelAllOffersBySellerInput) Hash() (common.Hash, error) {
	uint256Ty, _ := abi.NewType("int256", "int256", nil)

	arguments := abi.Arguments{
		{Type: uint256Ty},
	}

	playerId, _ := new(big.Int).SetString(string(b.PlayerId), 10)
	if playerId == nil {
		return common.Hash{}, errors.New("invalid playerId")
	}

	bytes, err := arguments.Pack(
		playerId,
	)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(bytes), nil
}

func (b CancelAllOffersBySellerInput) SignerAddress() (common.Address, error) {
	hash, err := b.Hash()
	if err != nil {
		return common.Address{}, err
	}
	sign, err := hex.DecodeString(b.Signature)
	if err != nil {
		return common.Address{}, err
	}
	return helper.AddressFromHashAndSignature(hash, sign)
}
