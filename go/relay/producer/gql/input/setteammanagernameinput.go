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

type SetTeamManagerNameInput struct {
	Signature string
	TeamId    graphql.ID
	Name      string
}

func (b SetTeamManagerNameInput) Hash() (common.Hash, error) {
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	stringTy, _ := abi.NewType("string", "string", nil)

	arguments := abi.Arguments{
		{Type: uint256Ty},
		{Type: stringTy},
	}

	teamId, _ := new(big.Int).SetString(string(b.TeamId), 10)
	if teamId == nil {
		return common.Hash{}, errors.New("Invalid TeamId")
	}

	bytes, err := arguments.Pack(
		teamId,
		b.Name,
	)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(bytes), nil
	return common.Hash{}, nil
}

func (b SetTeamManagerNameInput) IsValidSignature() (bool, error) {
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

func (b SetTeamManagerNameInput) SignerAddress() (common.Address, error) {
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
