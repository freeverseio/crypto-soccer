package assets

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
	Backend          *backends.SimulatedBackend
	Ops              *bind.TransactOpts
	PrivateKey       *ecdsa.PrivateKey
	playerStatesAddr common.Address
	assets           *Assets
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
	return &SimulatedBlockchain{blockchain, auth, genesisPrivateKey, common.Address{}, nil}
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
func (blockchain *SimulatedBlockchain) DeployAssets(owner *ecdsa.PrivateKey) {
	blockchain.playerStatesAddr = common.HexToAddress("0x83a909262608c650bd9b0ae06e29d90d0f67ac5e")
	//Deploy contract
	_, _, contract, err := DeployAssets(
		bind.NewKeyedTransactor(owner),
		blockchain.Backend,
		blockchain.playerStatesAddr,
	)
	AssertNoErr(err, "DeployAssets failed")
	blockchain.Commit()
	blockchain.assets = contract
}
func (blockchain *SimulatedBlockchain) CreateTeam(name string, from *ecdsa.PrivateKey) {
	auth := bind.NewKeyedTransactor(from)
	_, err := blockchain.assets.CreateTeam(
		&bind.TransactOpts{
			From:   auth.From,
			Signer: auth.Signer,
		},
		name,
		blockchain.playerStatesAddr)
	AssertNoErr(err, "Error creating Team")
	blockchain.Commit()
}
func (blockchain *SimulatedBlockchain) CountTeams() *big.Int {
	count, err := blockchain.assets.CountTeams(nil)
	AssertNoErr(err, "Error calling CountTeams")
	blockchain.Commit()
	return count
}
func (blockchain *SimulatedBlockchain) ScanTeamCreated() []AssetsTeamCreated {
	events, err := blockchain.assets.ScanTeamCreated()
	AssertNoErr(err, "Error calling ScanTeamCreatedTeam")
	blockchain.Commit()
	return events
}

func DefaultSimulatedBlockchain() *SimulatedBlockchain {
	genesisGasLimit := uint64(8000029)
	genesisBalance := "1000000000000000000000" // 1000 eth in wei
	blockchain := NewSimulatedBlockchain(genesisGasLimit, genesisBalance)
	bob := blockchain.CreateAccountWithBalance("1000000000000000000") // 1 eth
	blockchain.DeployAssets(bob)
	return blockchain
}
