package processor_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/market/notary/testutils"
)

func TestChangeOwnership(t *testing.T) {
	ganache := testutils.NewGanache()
	owner := ganache.CreateAccountWithBalance("1000000000000000000") // 1 eth
	ganache.DeployContracts(owner)

	alice := ganache.CreateAccountWithBalance("50000000000000000000") // 50 eth
	bob := ganache.CreateAccountWithBalance("50000000000000000000")   // 50 eth

	_, err := ganache.Assets.CreateTeam(
		bind.NewKeyedTransactor(owner),
		"Barca",
		crypto.PubkeyToAddress(alice.PublicKey))
	if err != nil {
		t.Fatal(err)
	}
	name, err := ganache.Assets.GetTeamName(nil, big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	if name != "Barca" {
		t.Errorf("Expected Barca got %v", name)
	}
	_, err = ganache.Assets.CreateTeam(
		bind.NewKeyedTransactor(owner),
		"Madrid",
		crypto.PubkeyToAddress(bob.PublicKey))
	if err != nil {
		t.Fatal(err)
	}
	name, err = ganache.Assets.GetTeamName(nil, big.NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	if name != "Madrid" {
		t.Errorf("Expected Madrid got %v", name)
	}

	_, err = ganache.Assets.TransferPlayer(
		bind.NewKeyedTransactor(alice),
		big.NewInt(1),
		big.NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
}

// func TestProcess(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	ganache := testutils.NewGanache()
// 	owner := ganache.CreateAccountWithBalance("1000000000000000000") // 1 eth
// 	ganache.DeployContracts(owner)

// 	processor, err := processor.NewProcessor(sto, ganache.Client, ganache.Assets, owner)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	alice := ganache.CreateAccountWithBalance("50000000000000000000") // 50 eth
// 	bob := ganache.CreateAccountWithBalance("50000000000000000000")   // 50 eth

// 	_, err = ganache.Assets.CreateTeam(
// 		bind.NewKeyedTransactor(owner),
// 		"Barca",
// 		crypto.PubkeyToAddress(alice.PublicKey))
// 	// ganache.CreateTeam("Barca", alice)
// 	ganache.CreateTeam("Madrid", bob)

// 	sto.CreateSellOrder(storage.SellOrder{1, 100})
// 	sto.CreateBuyOrder(storage.BuyOrder{1, 100, 2})
// 	processor.Process()
// }
