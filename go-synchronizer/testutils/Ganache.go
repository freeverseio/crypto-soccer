package testutils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"
)

type Ganache struct {
	Client  *ethclient.Client
	time    *Testutils
	Updates *updates.Updates
	Leagues *leagues.Leagues
	Engine  *engine.Engine
	Owner   *ecdsa.PrivateKey
}

func NewGanache() *Ganache {
	client, err := ethclient.Dial("http://localhost:8545")
	AssertNoErr(err, "Error connecting to ganache")
	creatorPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	AssertNoErr(err, "Failed converting private key to ECSDA")

	// deploy time.go to fake block transactions and be able to advance blocknumber in ganache
	_, _, time, err := DeployTestutils(
		bind.NewKeyedTransactor(creatorPrivateKey),
		client,
	)
	AssertNoErr(err, "DeployTime failed")

	return &Ganache{
		client,
		time,
		nil,
		nil,
		nil,
		creatorPrivateKey,
	}
}
func (ganache *Ganache) Advance(blockCount int) {
	for i := 0; i < blockCount; i++ {
		_, err := ganache.time.Increase(bind.NewKeyedTransactor(ganache.Owner))
		AssertNoErr(err, "Error in Advance()")
	}
}
func (ganache *Ganache) CreateAccountWithBalance(wei string) *ecdsa.PrivateKey {
	value := new(big.Int)
	value.SetString(wei, 10)
	privateKey, err := crypto.GenerateKey()
	AssertNoErr(err, "Failed generating key")
	toAddress := ganache.Public(privateKey)
	_, err = ganache.TransferWei(value, ganache.Owner, toAddress)
	AssertNoErr(err, "Failed transferring wei")

	return privateKey
}
func (ganache *Ganache) Public(addr *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(addr.PublicKey)
}
func (ganache *Ganache) GetNonce(from *ecdsa.PrivateKey) uint64 {
	publicKey := from.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	nonce, err := ganache.Client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(*publicKeyECDSA))
	AssertNoErr(err, "Failed obtaining pending nonce")
	return nonce
}
func (ganache *Ganache) TransferWei(wei *big.Int, from *ecdsa.PrivateKey, to common.Address) (*types.Transaction, error) {
	nonce := ganache.GetNonce(from)
	gasLimit := uint64(21000)
	gasPrice, err := ganache.Client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("TransferWei: Failed obtaining suggested gas price, using 2GWei")
		gasPrice := new(big.Int)
		gasPrice.SetString("2000000000", 10)
	}
	var data []byte
	tx := types.NewTransaction(nonce, to, wei, gasLimit, gasPrice, data)
	chainID, err := ganache.Client.NetworkID(context.Background())
	AssertNoErr(err, "TransferWei: Failed obtaining chainID")
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), from)
	err = ganache.Client.SendTransaction(context.Background(), signedTx)
	return signedTx, err
}
func (ganache *Ganache) GetLastBlockNumber() int64 {
	header, err := ganache.Client.HeaderByNumber(context.Background(), nil)
	AssertNoErr(err, "Failed GetLastBlockNumber")
	return header.Number.Int64()
}
func (ganache *Ganache) GetBalance(address common.Address) *big.Int {
	lastBlock := big.NewInt(ganache.GetLastBlockNumber())
	balance, err := ganache.Client.BalanceAt(context.Background(), address, lastBlock)
	AssertNoErr(err, "Failed GetBalance")
	return balance
}

func (ganache *Ganache) DeployContracts(owner *ecdsa.PrivateKey) {
	leaguesAddress, _, leaguesContract, err := leagues.DeployLeagues(
		bind.NewKeyedTransactor(owner),
		ganache.Client,
	)
	AssertNoErr(err, "DeployLeagues failed")
	fmt.Println("Leagues deployed at:", leaguesAddress.Hex())

	updatesAddress, _, updatesContract, err := updates.DeployUpdates(
		bind.NewKeyedTransactor(owner),
		ganache.Client,
	)
	AssertNoErr(err, "DeployUpdates failed")
	fmt.Println("Updates deployed at:", updatesAddress.Hex())

	engineAddress, _, engineContract, err := engine.DeployEngine(
		bind.NewKeyedTransactor(owner),
		ganache.Client,
	)
	AssertNoErr(err, "DeployEngine failed")
	fmt.Println("Engine deployed at:", engineAddress.Hex())

	// setup
	_, err = leaguesContract.SetEngineAdress(bind.NewKeyedTransactor(owner), engineAddress)
	AssertNoErr(err, "Error setting engine contract in league contract")
	_, err = updatesContract.InitUpdates(bind.NewKeyedTransactor(owner), leaguesAddress)
	AssertNoErr(err, "Updates::InitUpdates(leagues) failed")

	ganache.Updates = updatesContract
	ganache.Leagues = leaguesContract
	ganache.Engine = engineContract
}

func (ganache *Ganache) Init() {
	// Initing
	_, err := ganache.Leagues.Init(bind.NewKeyedTransactor(ganache.Owner))
	AssertNoErr(err, "Error initializing leagues contract")
}

func (ganache *Ganache) InitOneTimezone() {
	// Initing
	_, err := ganache.Leagues.InitSingleTZ(bind.NewKeyedTransactor(ganache.Owner), 1)
	AssertNoErr(err, "Error initializing leagues contract")
}

//func (ganache *Ganache) CountTeams() *big.Int {
//	count, err := ganache.Assets.CountTeams(nil)
//	AssertNoErr(err, "Error calling CountTeams")
//	return count
//}

//func (ganache *Ganache) CountLeagues() *big.Int {
//	count, err := ganache.Leagues.LeaguesCount(nil)
//	AssertNoErr(err)
//	return count
//}
//func PrintTeamCreated(event assets.AssetsTeamCreated, ganache *Ganache) {
//	fmt.Println("team id:", event.Id.Int64(), "players: ", ganache.GetVirtualPlayers(event.Id))
//}
