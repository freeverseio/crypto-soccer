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

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/states"
)

type Ganache struct {
	Client        *ethclient.Client
	time          *Testutils
	statesAddress common.Address
	engineAddress common.Address
	Assets        *assets.Assets
	States        *states.States
	Leagues       *leagues.Leagues
	Owner         *ecdsa.PrivateKey
}

func NewGanache() *Ganache {
	client, err := ethclient.Dial("http://localhost:8545")
	AssertNoErr(err, "Error connecting to ganache")
	creatorPrivateKey, err := crypto.HexToECDSA("f1b3f8e0d52caec13491368449ab8d90f3d222a3e485aa7f02591bbceb5efba5")
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
		common.Address{},
		common.Address{},
		nil,
		nil,
		nil,
		creatorPrivateKey,
	}
}
func (ganache *Ganache) Advance(blockCount int) {
	for i := 0; i < blockCount; i++ {
		auth := bind.NewKeyedTransactor(ganache.Owner)
		_, err := ganache.time.Increase(
			&bind.TransactOpts{
				From:   auth.From,
				Signer: auth.Signer,
				//GasLimit: uint64(2000000000),
			})
		AssertNoErr(err, "Error in Advance()")
	}
}
func (ganache *Ganache) CreateAccountWithBalance(wei string) *ecdsa.PrivateKey {
	value := new(big.Int)
	value.SetString(wei, 10)
	privateKey, err := crypto.GenerateKey()
	AssertNoErr(err, "Failed generating key")
	toAddress := CommonAddressFromPrivateKey(privateKey)
	_, err = ganache.TransferWei(value, ganache.Owner, toAddress)
	AssertNoErr(err, "Failed transferring wei")

	return privateKey
}
func (ganache *Ganache) Public(addr *ecdsa.PrivateKey) common.Address {
	publicKey := addr.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	return crypto.PubkeyToAddress(*publicKeyECDSA)
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
func (ganache *Ganache) deployStates(owner *ecdsa.PrivateKey) {
	address, _, contract, err := states.DeployStates(
		bind.NewKeyedTransactor(owner),
		ganache.Client,
	)
	AssertNoErr(err, "DeployStates failed")
	ganache.States = contract
	ganache.statesAddress = address
}
func (ganache *Ganache) deployEngine(owner *ecdsa.PrivateKey) {
	address, _, contract, err := engine.DeployEngine(
		bind.NewKeyedTransactor(owner),
		ganache.Client,
	)
	AssertNoErr(err, "DeployEngine failed")
	_ = contract
	ganache.engineAddress = address
}
func (ganache *Ganache) deployAssets(owner *ecdsa.PrivateKey) {
	_, _, contract, err := assets.DeployAssets(
		bind.NewKeyedTransactor(owner),
		ganache.Client,
		ganache.statesAddress,
	)
	AssertNoErr(err, "DeployAssets failed")
	ganache.Assets = contract
}
func (ganache *Ganache) deployLeagues(owner *ecdsa.PrivateKey) {
	_, _, contract, err := leagues.DeployLeagues(
		bind.NewKeyedTransactor(owner),
		ganache.Client,
		ganache.engineAddress,
		ganache.statesAddress,
	)
	AssertNoErr(err, "DeployStates failed")
	ganache.Leagues = contract
}
func (ganache *Ganache) DeployContracts(owner *ecdsa.PrivateKey) {
	ganache.deployStates(owner)
	ganache.deployEngine(owner)
	ganache.deployAssets(owner)
	ganache.deployLeagues(owner)
}
func (ganache *Ganache) CreateTeam(name string, from *ecdsa.PrivateKey) {
	auth := bind.NewKeyedTransactor(from)
	_, err := ganache.Assets.CreateTeam(
		&bind.TransactOpts{
			From:   auth.From,
			Signer: auth.Signer,
			//GasLimit: uint64(2000000000),
		},
		name,
		ganache.statesAddress)
	AssertNoErr(err, "Error creating Team ", name)
}
func (ganache *Ganache) getVirtualPlayerId(teamId *big.Int, posInTeam uint8) int64 {
	playerId, err := ganache.Assets.GenerateVirtualPlayerId(
		&bind.CallOpts{},
		teamId,
		posInTeam,
	)
	AssertNoErr(err, "Error getting virtual player id in pos ", posInTeam, " for team ", teamId)
	return playerId.Int64()
}
func (ganache *Ganache) getVirtualPlayerState(playerId int64) *big.Int {
	playerState, err := ganache.Assets.GenerateVirtualPlayerState(
		&bind.CallOpts{},
		big.NewInt(playerId),
	)
	AssertNoErr(err, "Error getting virtual player state for id ", playerId)
	return playerState
}
func (ganache *Ganache) GetVirtualPlayers(teamId *big.Int) (players map[int64]*big.Int) {
	players = make(map[int64]*big.Int)
	for i := 0; i < 11; i++ {
		playerId := ganache.getVirtualPlayerId(teamId, uint8(i))
		playerState := ganache.getVirtualPlayerState(playerId)
		players[playerId] = playerState
	}
	return
}
func (ganache *Ganache) CreateLeague(teamIds []int64, from *ecdsa.PrivateKey) {
	leagueId := ganache.CountLeagues()
	initBlock := big.NewInt(ganache.GetLastBlockNumber())
	step := big.NewInt(1)
	var tactics [][3]uint8 // {[3]uint8{4, 4, 2}, [3]uint8{4, 3, 3}},
	var teamIdsBig []*big.Int
	for _, teamId := range teamIds {
		teamIdsBig = append(teamIdsBig, big.NewInt(teamId))
		if teamId%3 == 0 {
			tactics = append(tactics, [3]uint8{3, 4, 3})
		} else if teamId%2 == 0 {
			tactics = append(tactics, [3]uint8{4, 3, 3})
		} else {
			tactics = append(tactics, [3]uint8{4, 4, 2})
		}
	}
	auth := bind.NewKeyedTransactor(from)
	tx, err := ganache.Leagues.Create(
		&bind.TransactOpts{
			From:   auth.From,
			Signer: auth.Signer,
			//GasLimit: uint64(2000000000),
		}, leagueId, initBlock, step, teamIdsBig, tactics)
	_ = tx
	AssertNoErr(err)
}
func (ganache *Ganache) CountTeams() *big.Int {
	count, err := ganache.Assets.CountTeams(nil)
	AssertNoErr(err, "Error calling CountTeams")
	return count
}
func (ganache *Ganache) CountLeagues() *big.Int {
	count, err := ganache.Leagues.LeaguesCount(nil)
	AssertNoErr(err)
	return count
}
func PrintTeamCreated(event assets.AssetsTeamCreated, ganache *Ganache) {
	fmt.Println("team id:", event.Id.Int64(), "players: ", ganache.GetVirtualPlayers(event.Id))
}
