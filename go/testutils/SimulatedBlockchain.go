package testutils

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"
)

// AssertNoErr - log fatal and panic on error and print params
func AssertNoErr(err error, params ...interface{}) {
	if err != nil {
		log.Fatal(err, params)
	}
}

func CommonAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) common.Address {
	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey) // type casting to ecdsa
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

type SimulatedBlockchain struct {
	Backend    *backends.SimulatedBackend
	Ops        *bind.TransactOpts
	PrivateKey *ecdsa.PrivateKey
	Updates    *updates.Updates
	Leagues    *leagues.Leagues
}

func NewSimulatedBlockchain(genesisGasLimit uint64, genesisBalanceWei string) *SimulatedBlockchain {
	genesisPrivateKey, err := crypto.GenerateKey()
	genesisBalance := new(big.Int)
	genesisBalance.SetString(genesisBalanceWei, 10)
	AssertNoErr(err, "failed generating key")
	auth := bind.NewKeyedTransactor(genesisPrivateKey)
	blockchain := backends.NewSimulatedBackend(
		map[common.Address]core.GenesisAccount{auth.From: {Balance: genesisBalance}},
		genesisGasLimit,
	)

	return &SimulatedBlockchain{
		blockchain,
		auth,
		genesisPrivateKey,
		nil,
		nil,
	}
}
func (blockchain *SimulatedBlockchain) Commit() {
	blockchain.Backend.Commit()
}

func (blockchain *SimulatedBlockchain) CreateAccountWithBalance(wei string) *ecdsa.PrivateKey {
	value := new(big.Int)
	value.SetString(wei, 10)
	gasLimit := uint64(21000)
	gasPrice, err := blockchain.Backend.SuggestGasPrice(context.Background())
	AssertNoErr(err)
	var data []byte

	privateKey, err := crypto.GenerateKey()
	AssertNoErr(err, "failed generating key")
	toAddress := CommonAddressFromPrivateKey(privateKey)

	nonce, err := blockchain.Backend.PendingNonceAt(context.Background(), blockchain.Ops.From)
	AssertNoErr(err)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	//chainID := big.NewInt(1337)
	//signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), genesisPrivateKey)
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, blockchain.PrivateKey)
	AssertNoErr(err)
	err = blockchain.Backend.SendTransaction(context.Background(), signedTx)
	AssertNoErr(err)
	blockchain.Commit()
	return privateKey
}

func (blockchain *SimulatedBlockchain) deployEngine(owner *ecdsa.PrivateKey) (common.Address, *engine.Engine) {
	address, _, contract, err := engine.DeployEngine(
		bind.NewKeyedTransactor(owner),
		blockchain.Backend,
	)
	AssertNoErr(err, "DeployEngine failed")
	blockchain.Commit()
	return address, contract
}

func (blockchain *SimulatedBlockchain) deployLeagues(owner *ecdsa.PrivateKey) (common.Address, *leagues.Leagues) {
	address, _, contract, err := leagues.DeployLeagues(
		bind.NewKeyedTransactor(owner),
		blockchain.Backend,
	)
	AssertNoErr(err, "DeployAssets failed")
	blockchain.Commit()
	return address, contract
}

func (blockchain *SimulatedBlockchain) deployUpdates(owner *ecdsa.PrivateKey) (common.Address, *updates.Updates) {
	address, _, contract, err := updates.DeployUpdates(
		bind.NewKeyedTransactor(owner),
		blockchain.Backend,
	)
	AssertNoErr(err, "DeployUpdates failed")
	blockchain.Commit()
	return address, contract
}

func (blockchain *SimulatedBlockchain) DeployContracts(owner *ecdsa.PrivateKey) {
	engineAddress, _ := blockchain.deployEngine(owner)
	leaguesAddress, leaguesContract := blockchain.deployLeagues(owner)
	_, updatesContract := blockchain.deployUpdates(owner)

	_, err := leaguesContract.SetEngineAdress(bind.NewKeyedTransactor(owner), engineAddress)
	AssertNoErr(err, "Error setting engine contract in league contract")
	_, err = leaguesContract.InitSingleTZ(bind.NewKeyedTransactor(owner), 1)
	AssertNoErr(err, "Error initializing leagues contract")
	_, err = updatesContract.InitUpdates(bind.NewKeyedTransactor(owner), leaguesAddress)
	AssertNoErr(err, "Updates::InitUpdates(leagues) failed")

	blockchain.Updates = updatesContract
	blockchain.Leagues = leaguesContract
}

//func (blockchain *SimulatedBlockchain) CreateTeam(name string, from *ecdsa.PrivateKey) {
//	auth := bind.NewKeyedTransactor(from)
//	_, err := blockchain.Assets.CreateTeam(
//		&bind.TransactOpts{
//			From:   auth.From,
//			Signer: auth.Signer,
//		},
//		name,
//		blockchain.statesAddress)
//	AssertNoErr(err, "Error creating Team")
//	blockchain.Commit()
//}
//func (blockchain *SimulatedBlockchain) CountTeams() *big.Int {
//	count, err := blockchain.Assets.CountTeams(nil)
//	AssertNoErr(err, "Error calling CountTeams")
//	blockchain.Commit()
//	return count
//}
func DefaultSimulatedBlockchain() *SimulatedBlockchain {
	genesisGasLimit := uint64(8000029)
	genesisBalance := "1000000000000000000000" // 1000 eth in wei
	blockchain := NewSimulatedBlockchain(genesisGasLimit, genesisBalance)
	owner := blockchain.CreateAccountWithBalance("1000000000000000000") // 1 eth
	blockchain.DeployContracts(owner)
	return blockchain
}
