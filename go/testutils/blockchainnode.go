package testutils

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/contracts/truffle"
	"github.com/freeverseio/crypto-soccer/go/helper"

	log "github.com/sirupsen/logrus"
)

type BlockchainNode struct {
	Client    *ethclient.Client
	Owner     *ecdsa.PrivateKey
	Contracts *contracts.Contracts
}

// AssertNoErr - log fatal and panic on error and print params
func AssertNoErr(err error, params ...interface{}) {
	if err != nil {
		log.Fatal(err, params)
	}
}

func NewBlockchainNodeDeployAndInitAt(address string) (*BlockchainNode, error) {
	node, err := NewBlockchainNodeAt(address)
	if err != nil {
		return nil, err
	}
	err = node.DeployContracts(node.Owner)
	if err != nil {
		return nil, err
	}
	err = node.InitOneTimezone(1)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func NewBlockchainNodeAt(address string) (*BlockchainNode, error) {
	client, err := ethclient.Dial(address)
	if err != nil {
		return nil, err
	}
	creatorPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	if err != nil {
		return nil, err
	}

	return &BlockchainNode{
		client,
		creatorPrivateKey,
		nil,
	}, nil
}

func NewBlockchainNodeDeployAndInit() (*BlockchainNode, error) {
	return NewBlockchainNodeDeployAndInitAt("http://localhost:8545")
}

func NewBlockchainNode() (*BlockchainNode, error) {
	return NewBlockchainNodeAt("http://localhost:8545")
}

func (b *BlockchainNode) DeployContracts(owner *ecdsa.PrivateKey) error {
	directoryAddress, err := truffle.DeplyByTruffle()
	if err != nil {
		return err
	}

	b.Contracts, err = contracts.NewByDirectoryAddress(
		b.Client,
		directoryAddress,
	)

	return err
}

func (b *BlockchainNode) Init() error {
	// Initing
	tx, err := b.Contracts.Assets.InitTZs(bind.NewKeyedTransactor(b.Owner))
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx, 10)
	if err != nil {
		return err
	}
	return nil
}

func (b *BlockchainNode) InitOneTimezone(timezoneIdx uint8) error {
	// Initing
	tx, err := b.Contracts.Assets.InitSingleTZ(bind.NewKeyedTransactor(b.Owner), timezoneIdx)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx, 10)
	if err != nil {
		return err
	}
	return nil
}
