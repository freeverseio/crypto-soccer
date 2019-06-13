package process

import (
	"math/big"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Assets interface {
	CountTeams(opts *bind.CallOpts) (*big.Int, error)
}