package input

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type GeneratePlayerIdsInput struct {
	Signature string
	Seed      int32
}

func (b GeneratePlayerIdsInput) Hash() (common.Hash, error) {
	int32Ty, _ := abi.NewType("int32", "int32", nil)
	arguments := abi.Arguments{
		{
			Type: int32Ty,
		},
	}

	bytes, err := arguments.Pack(b.Seed)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(bytes), nil
	// copy(hash[:], crypto.Keccak256Hash(bytes).Bytes())
	// ss := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(hash), hash)
	// return crypto.Keccak256Hash([]byte(ss)), nil
}
