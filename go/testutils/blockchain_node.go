package testutils

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/crypto-soccer/go/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go/contracts/engineprecomp"
	"github.com/freeverseio/crypto-soccer/go/contracts/evolution"
	"github.com/freeverseio/crypto-soccer/go/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go/contracts/market"
	"github.com/freeverseio/crypto-soccer/go/contracts/updates"
	"github.com/freeverseio/crypto-soccer/go/helper"
)

type BlockchainNode struct {
	Client        *ethclient.Client
	Assets        *assets.Assets
	Updates       *updates.Updates
	Leagues       *leagues.Leagues
	Engine        *engine.Engine
	EnginePreComp *engineprecomp.Engineprecomp
	Market        *market.Market
	Evolution     *evolution.Evolution
	Owner         *ecdsa.PrivateKey
}

// AssertNoErr - log fatal and panic on error and print params
func AssertNoErr(err error, params ...interface{}) {
	if err != nil {
		log.Fatal(err, params)
	}
}

func NewBlockchainNodeDeployAndInit() (*BlockchainNode, error) {
	node, err := NewBlockchainNode()
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
		nil,
		nil,
		nil,
		creatorPrivateKey,
	}, nil
}

func (b *BlockchainNode) DeployContracts(owner *ecdsa.PrivateKey) error {
	assetsAddress, tx10, assetsContract, err := assets.DeployAssets(
		bind.NewKeyedTransactor(owner),
		b.Client,
	)
	AssertNoErr(err, "DeployAssets failed")
	fmt.Println("Assets deployed at:", assetsAddress.Hex())
	if err != nil {
		return err
	}

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

	evolutionAddress, tx30, evolutionContract, err := evolution.DeployEvolution(
		bind.NewKeyedTransactor(owner),
		b.Client,
	)
	AssertNoErr(err, "DeployEvolution failed")
	fmt.Println("Evolution deployed at:", evolutionAddress.Hex())
	if err != nil {
		return err
	}

	engineprecompAddress, tx31, enginePreComp, err := engineprecomp.DeployEngineprecomp(
		bind.NewKeyedTransactor(owner),
		b.Client,
	)
	AssertNoErr(err, "DeployEngineprecomp failed")
	fmt.Println("Engineprecomp deployed at:", engineprecompAddress.Hex())
	if err != nil {
		return err
	}

	_, err = helper.WaitReceipt(b.Client, tx10, 10)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx0, 10)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx1, 10)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx2, 10)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx3, 10)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx30, 10)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx31, 10)
	if err != nil {
		return err
	}
	// setup
	tx0, err = leaguesContract.SetEngineAdress(bind.NewKeyedTransactor(owner), engineAddress)
	AssertNoErr(err, "Error setting engine contract in league contract")
	tx2, err = marketContract.SetAssetsAddress(bind.NewKeyedTransactor(owner), assetsAddress)
	AssertNoErr(err, "Error setting Assets address to market")
	tx1, err = updatesContract.InitUpdates(bind.NewKeyedTransactor(owner), assetsAddress)
	AssertNoErr(err, "Updates::InitUpdates(leagues) failed")
	tx3, err = evolutionContract.SetEngine(bind.NewKeyedTransactor(owner), engineAddress)
	AssertNoErr(err, "Error setting engine contract in evolution contract")
	tx30, err = engineContract.SetPreCompAddr(bind.NewKeyedTransactor(owner), engineprecompAddress)
	AssertNoErr(err, "Error setting engineprecomp contract in engine contract")

	_, err = helper.WaitReceipt(b.Client, tx0, 10)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx1, 10)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx2, 10)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx3, 10)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx30, 10)
	if err != nil {
		return err
	}

	b.Updates = updatesContract
	b.Leagues = leaguesContract
	b.Engine = engineContract
	b.Market = marketContract
	b.Assets = assetsContract
	b.Evolution = evolutionContract
	b.EnginePreComp = enginePreComp

	return nil
}

func (b *BlockchainNode) Init() error {
	// Initing
	tx, err := b.Assets.Init(bind.NewKeyedTransactor(b.Owner))
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
	tx, err := b.Assets.InitSingleTZ(bind.NewKeyedTransactor(b.Owner), timezoneIdx)
	if err != nil {
		return err
	}
	_, err = helper.WaitReceipt(b.Client, tx, 10)
	if err != nil {
		return err
	}
	return nil
}
