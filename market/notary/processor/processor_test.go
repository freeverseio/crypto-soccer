package processor_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/market/notary/processor"
	"github.com/freeverseio/crypto-soccer/market/notary/storage"
	"github.com/freeverseio/crypto-soccer/market/notary/testutils"
)

func TestChangeOwnership(t *testing.T) {
	ganache := testutils.NewGanache()
	alice := ganache.Alice
	bob := ganache.Bob

	ganache.CreateTeam("Barca", alice)
	ganache.CreateTeam("Madrid", bob)

	var player = big.NewInt(1)
	originOwner := ganache.GetPlayerOwner(player)
	if originOwner != ganache.Public(alice) {
		t.Fatalf("Expected owner ALICE but got %v", originOwner)
	}
	err := ganache.TransferPlayer(player, big.NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	targetOwner := ganache.GetPlayerOwner(player)
	if targetOwner != ganache.Public(bob) {
		t.Fatalf("Expectedf owner BOB but got %v", targetOwner)
	}
}

func TestProcess(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache := testutils.NewGanache()
	owner := ganache.CreateAccountWithBalance("1000000000000000000") // 1 eth
	ganache.DeployContracts(owner)

	processor, err := processor.NewProcessor(sto, ganache.Client, ganache.Assets, owner)
	if err != nil {
		t.Fatal(err)
	}

	alice := ganache.CreateAccountWithBalance("50000000000000000000") // 50 eth
	bob := ganache.CreateAccountWithBalance("50000000000000000000")   // 50 eth

	ganache.CreateTeam("Barca", alice)
	ganache.CreateTeam("Madrid", bob)

	var player = big.NewInt(1)
	originOwner, err := ganache.Assets.GetPlayerOwner(&bind.CallOpts{}, player)
	if err != nil {
		t.Fatal(err)
	}
	if originOwner != crypto.PubkeyToAddress(alice.PublicKey) {
		t.Fatalf("Expectedf originOwner ALICE but got %v", originOwner)
	}
	sto.CreateSellOrder(storage.SellOrder{1, 100})
	processor.Process()
	targetOwner, err := ganache.Assets.GetPlayerOwner(&bind.CallOpts{}, player)
	if err != nil {
		t.Fatal(err)
	}
	if targetOwner != crypto.PubkeyToAddress(alice.PublicKey) {
		t.Fatalf("Expectedf originOwner ALICE but got %v", targetOwner)
	}

	sto.CreateBuyOrder(storage.BuyOrder{1, 100, 2})
	processor.Process()
	targetOwner, err = ganache.Assets.GetPlayerOwner(&bind.CallOpts{}, player)
	if err != nil {
		t.Fatal(err)
	}
	if targetOwner != crypto.PubkeyToAddress(bob.PublicKey) {
		t.Fatalf("Expectedf originOwner BOB but got %v", targetOwner)
	}

	buyOrders, err := sto.GetBuyOrders()
	if err != nil {
		t.Fatal(err)
	}
	if len(buyOrders) != 0 {
		t.Fatalf("Expercted 0 but got %v", len(buyOrders))
	}
	sellOrders, err := sto.GetSellOrders()
	if err != nil {
		t.Fatal(err)
	}
	if len(sellOrders) != 0 {
		t.Fatalf("Expercted 0 but got %v", len(sellOrders))
	}
}
