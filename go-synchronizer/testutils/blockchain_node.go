package testutils

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/market"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"
)

type BlockchainNode struct {
	Client  *ethclient.Client
	Updates *updates.Updates
	Leagues *leagues.Leagues
	Engine  *engine.Engine
	Market  *market.Market
	Owner   *ecdsa.PrivateKey
}

func NewBlockchainNode() (*BlockchainNode, error) {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		return nil, err
	}
	creatorPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	if err != nil {
		return nil, err
	}

	return &BlockchainNode{
		client,
		nil,
		nil,
		nil,
		nil,
		creatorPrivateKey,
	}, nil
}

func (b *BlockchainNode) WaitReceipt(tx *types.Transaction, timeoutSec uint8) error {
	receiptTimeout := time.Second * time.Duration(timeoutSec)
	start := time.Now()
	ctx := context.TODO()
	var receipt *types.Receipt

	for receipt == nil && time.Now().Sub(start) < receiptTimeout {
		receipt, err := b.Client.TransactionReceipt(ctx, tx.Hash())
		if err == nil && receipt != nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return errors.New("Timeout waiting for receipt")
}

func (b *BlockchainNode) WaitReceipts(txs []*types.Transaction, timeoutSec uint8) error {
	for _, tx := range txs {
		err := b.WaitReceipt(tx, timeoutSec)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BlockchainNode) DeployContracts(owner *ecdsa.PrivateKey) error {
	leaguesAddress, tx0, leaguesContract, err := leagues.DeployLeagues(
		bind.NewKeyedTransactor(owner),
		b.Client,
	)
	AssertNoErr(err, "DeployLeagues failed")
	fmt.Println("Leagues deployed at:", leaguesAddress.Hex())
	if err != nil {
		return err
	}

	updatesAddress, tx1, updatesContract, err := updates.DeployUpdates(
		bind.NewKeyedTransactor(owner),
		b.Client,
	)
	AssertNoErr(err, "DeployUpdates failed")
	fmt.Println("Updates deployed at:", updatesAddress.Hex())
	if err != nil {
		return err
	}

	engineAddress, tx2, engineContract, err := engine.DeployEngine(
		bind.NewKeyedTransactor(owner),
		b.Client,
	)
	AssertNoErr(err, "DeployEngine failed")
	fmt.Println("Engine deployed at:", engineAddress.Hex())
	if err != nil {
		return err
	}

	marketAddress, tx3, marketContract, err := market.DeployMarket(
		bind.NewKeyedTransactor(owner),
		b.Client,
	)
	AssertNoErr(err, "DeployMarket failed")
	fmt.Println("Market deployed at:", marketAddress.Hex())
	if err != nil {
		return err
	}

	err = b.WaitReceipt(tx0, 10)
	if err != nil {
		return err
	}
	err = b.WaitReceipt(tx1, 10)
	if err != nil {
		return err
	}
	err = b.WaitReceipt(tx2, 10)
	if err != nil {
		return err
	}
	err = b.WaitReceipt(tx3, 10)
	if err != nil {
		return err
	}
	// setup
	tx0, err = leaguesContract.SetEngineAdress(bind.NewKeyedTransactor(owner), engineAddress)
	AssertNoErr(err, "Error setting engine contract in league contract")
	tx2, err = marketContract.SetAssetsAddress(bind.NewKeyedTransactor(owner), leaguesAddress)
	AssertNoErr(err, "Error setting Assets address to market")
	tx1, err = updatesContract.InitUpdates(bind.NewKeyedTransactor(owner), leaguesAddress)
	AssertNoErr(err, "Updates::InitUpdates(leagues) failed")

	err = b.WaitReceipt(tx0, 10)
	if err != nil {
		return err
	}
	err = b.WaitReceipt(tx1, 10)
	if err != nil {
		return err
	}
	err = b.WaitReceipt(tx2, 10)
	if err != nil {
		return err
	}

	b.Updates = updatesContract
	b.Leagues = leaguesContract
	b.Engine = engineContract
	b.Market = marketContract

	return nil
}

func (b *BlockchainNode) Init() error {
	// Initing
	tx, err := b.Leagues.Init(bind.NewKeyedTransactor(b.Owner))
	if err != nil {
		return err
	}
	err = b.WaitReceipt(tx, 10)
	if err != nil {
		return err
	}
	return nil
}

func (b *BlockchainNode) InitOneTimezone(timezoneIdx uint8) error {
	// Initing
	tx, err := b.Leagues.InitSingleTZ(bind.NewKeyedTransactor(b.Owner), timezoneIdx)
	if err != nil {
		return err
	}
	err = b.WaitReceipt(tx, 10)
	if err != nil {
		return err
	}
	return nil
}
