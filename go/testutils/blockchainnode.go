package testutils

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/contracts/truffle"
)

type BlockchainNode struct {
	Client    *ethclient.Client
	Owner     *ecdsa.PrivateKey
	Contracts *contracts.Contracts
}

func NewBlockchain() (*BlockchainNode, error) {
	bc := BlockchainNode{}
	var err error
	bc.Owner, err = crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	if err != nil {
		return nil, err
	}

	bc.Contracts, err = truffle.New()
	if err != nil {
		return nil, err
	}

	bc.Client = bc.Contracts.Client

	return &bc, nil
}
