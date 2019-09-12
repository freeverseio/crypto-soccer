package processor_test

import (
	"math/big"
	"testing"

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
	owner := ganache.GetPlayerOwner(player)
	if owner != ganache.Public(alice) {
		t.Fatalf("Expected owner ALICE but got %v", owner)
	}
	err := ganache.TransferPlayer(player, big.NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	owner = ganache.GetPlayerOwner(player)
	if owner != ganache.Public(bob) {
		t.Fatalf("Expectedf owner BOB but got %v", owner)
	}
	ganache.CreateTeam("Venice", alice)
	err = ganache.TransferPlayer(player, big.NewInt(3))
	if err != nil {
		t.Fatal(err)
	}
	owner = ganache.GetPlayerOwner(player)
	if owner != ganache.Public(alice) {
		t.Fatalf("Expected owner ALICE but got %v", owner)
	}
}

func TestProcess(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache := testutils.NewGanache()
	alice := ganache.Alice
	bob := ganache.Bob

	processor, err := processor.NewProcessor(sto, ganache.Client, ganache.Assets, ganache.Owner)
	if err != nil {
		t.Fatal(err)
	}

	ganache.CreateTeam("Barca", alice)
	ganache.CreateTeam("Madrid", bob)

	var player = big.NewInt(1)
	originOwner := ganache.GetPlayerOwner(player)
	if originOwner != ganache.Public(alice) {
		t.Fatalf("Expectedf originOwner ALICE but got %v", originOwner)
	}
	sto.CreateSellOrder(storage.SellOrder{
		PlayerId:   big.NewInt(1),
		Price:      big.NewInt(100),
		Rnd:        big.NewInt(4353),
		ValidUntil: big.NewInt(3),
		TypeOfTx:   3,
	})
	processor.Process()
	targetOwner := ganache.GetPlayerOwner(player)
	if targetOwner != crypto.PubkeyToAddress(alice.PublicKey) {
		t.Fatalf("Expectedf originOwner ALICE but got %v", targetOwner)
	}

	sto.CreateBuyOrder(storage.BuyOrder{
		PlayerId: big.NewInt(1),
		TeamId:   big.NewInt(2),
	})

	processor.Process()
	targetOwner = ganache.GetPlayerOwner(player)
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
