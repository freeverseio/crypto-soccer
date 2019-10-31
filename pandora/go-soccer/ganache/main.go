package main

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	lionel "github.com/freeverseio/go-soccer/ganache/lionel"
	stakers "github.com/freeverseio/go-soccer/ganache/stakers"
)

// CreateDir - dirs and subdirs if they don't exist
func CreateDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}

// AssertNoErr - log fatal and panic on error and print params
func AssertNoErr(err error, params ...interface{}) {
	if err != nil {
		log.Fatal(err, params)
		panic(err)
	}
}

// DeployLionel - deploys lionel contract
func DeployLionel(client *ethclient.Client, privateKey *ecdsa.PrivateKey) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	AssertNoErr(err)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	AssertNoErr(err)

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	//auth.GasLimit = uint64(30000000) // in units
	auth.GasPrice = gasPrice

	log.Println("deploying lionel...")
	address, tx, instance, err := lionel.DeployLionel(auth, client)
	AssertNoErr(err)

	stakersAddress, err := instance.Stakers(&bind.CallOpts{})
	AssertNoErr(err)
	log.Printf("Lionel address: %v\n", address.Hex())
	log.Printf("Stakers address: %v\n", stakersAddress.Hex())
	log.Printf("Transaction hash: %v\n", tx.Hash().Hex())

	// verify lionel address is the same when queried from Stakers contract
	stakersInstance, err := stakers.NewStakers(stakersAddress, client)
	AssertNoErr(err)

	gameAddress, err := stakersInstance.Game(&bind.CallOpts{})
	if gameAddress != address {
		log.Fatal("Lionel adresses differ. Expected: ", address.Hex(), " but goy ", gameAddress.Hex())
		panic(-1)
	}
}

func importAccounts(client *ethclient.Client, accounts []string) {
	for _, addr := range accounts {
		// create keystores
		walletPath := "/tmp/freeverse/test_wallets"
		walletPwd := "111111"
		CreateDir(walletPath)
		ks := keystore.NewKeyStore(walletPath, keystore.StandardScryptN, keystore.StandardScryptP)

		pKey, err := crypto.HexToECDSA(addr)
		AssertNoErr(err)
		account, err := ks.ImportECDSA(pKey, walletPwd)
		AssertNoErr(err)
		header, err := client.HeaderByNumber(context.Background(), nil)
		lastBlock := big.NewInt(header.Number.Int64())
		balance, err := client.BalanceAt(context.Background(), account.Address, lastBlock)
		log.Printf("account imported: %v balance: %v Wei\n", account.Address.Hex(), balance)
	}
}

func main() {
	ownerPrivateKey := "f1b3f8e0d52caec13491368449ab8d90f3d222a3e485aa7f02591bbceb5efba5"

	stakersPrivateKeys := []string{
		"91821f9af458d612362136648fc8552a47d8289c0f25a8a1bf0860510332cef9",
		"bb32062807c162a5243dc9bcf21d8114cb636c376596e1cf2895ec9e5e3e0a68",
		"95ce6122165d94aa51b0fcf51021895b39b0ff291aa640c803d5401bd87894d5",
		"3af93668029f95d526fc1d2bdefccc120bfe1d26a0462d268e8f6b2f71402ba3",
		"3b24a4fdf2e6e1375008c387c5456ce00cb0772435ae1938c2fe833103393b9a",
		"cba858feeb49e1ca8053a5213987a22c3ee83d9f9f396e138940a018dd837ebb",
		"df48bfda4cb4b4e094803e923836a8538fbf607da79f6e46d68cdd43fb2f3f88",
		"487efb6249a8a4d45a19c8e5d1e5c7d3f6610a7e69f8f81ddcf368f9a0c0d6d5",
		"bb4cce73db59f456ea427e5862fdb0d5bc038a7d0b930cbb45e1c4f6d122289e",
	}

	client, err := ethclient.Dial("http://localhost:8545")
	AssertNoErr(err)
	privateKey, err := crypto.HexToECDSA(ownerPrivateKey)
	AssertNoErr(err)
	DeployLionel(client, privateKey)
	log.Println("Importing accounts:")
	importAccounts(client, stakersPrivateKeys)
}
