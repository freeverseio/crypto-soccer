package input

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/graph-gophers/graphql-go"
)

type GetWorldPlayersInput struct {
	Signature string
	TeamId    graphql.ID
}

func (b GetWorldPlayersInput) Hash() (common.Hash, error) {
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	arguments := abi.Arguments{
		{
			Type: uint256Ty,
		},
	}

	teamId, _ := new(big.Int).SetString(string(b.TeamId), 10)
	if teamId == nil {
		return common.Hash{}, errors.New("Invalid TeamId")
	}

	bytes, err := arguments.Pack(teamId)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(bytes), nil
	// copy(hash[:], crypto.Keccak256Hash(bytes).Bytes())
	// ss := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	// return crypto.Keccak256Hash([]byte(ss)), nil
}
