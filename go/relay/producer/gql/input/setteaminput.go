package input

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/graph-gophers/graphql-go"
)

type SetTeamNameInput struct {
	Signature string
	TeamId    graphql.ID
	Name      string
}

func (b SetTeamNameInput) Hash() (common.Hash, error) {
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
