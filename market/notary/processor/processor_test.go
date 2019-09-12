package processor_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/market/notary/processor"
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

// func TestRSV(t *testing.T) {
// 	r, s, v := processor.RSV("0x405c83733f474f6919032fd41bd2e37b1a3be444bc52380c0e3f4c79ce8245ce229b4b0fe3a9798b5aad5f8df5c6acc07e4810f1a111d7712bf06aee7c7384001b")
// 	t.Fatalf("r = %v, s = %v, v = %v", r, s, v)
// }

func TestHashPRivateMessage(t *testing.T) {
	ganache := testutils.NewGanache()
	processor, err := processor.NewProcessor(nil, ganache.Client, ganache.Assets, ganache.Owner)
	if err != nil {
		t.Fatal(err)
	}
	currencyId := uint8(1)
	price := big.NewInt(41234)
	rnd := big.NewInt(42321)
	privHash, err := processor.HashPrivateMsg(
		currencyId,
		price,
		rnd,
	)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(privHash[:])
	if result != "4200de738160a9e6b8f69648fbb7feb323f73fac5acff1b7bb546bb7ac3591fa" {
		t.Fatalf("Hash private error %v", result)
	}
}

func TestBuildPufForSaleMessage(t *testing.T) {
	ganache := testutils.NewGanache()
	processor, err := processor.NewProcessor(nil, ganache.Client, ganache.Assets, ganache.Owner)
	if err != nil {
		t.Fatal(err)
	}
	validUntil := big.NewInt(2000000000)
	playerId := big.NewInt(10)
	typeOfTx := uint8(1)
	currencyId := uint8(1)
	price := big.NewInt(41234)
	rnd := big.NewInt(42321)

	privHash, err := processor.HashPrivateMsg(
		currencyId,
		price,
		rnd,
	)
	if err != nil {
		t.Fatal(err)
	}
	hash, err := processor.HashBuyerMessage(
		privHash,
		validUntil,
		playerId,
		typeOfTx,
	)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(hash[:])
	if result != "ff3497f25b47dbc25101237ad159a698f8fee96d1873b844dcac6d84a72b6dc0" {
		t.Fatalf("Hash error %v", result)
	}
}

// func TestProcess(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	ganache := testutils.NewGanache()
// 	alice, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	value := new(big.Int)
// 	value.SetString("50000000000000000000", 10)
// 	_, err = ganache.TransferWei(value, ganache.Owner, ganache.Public(alice))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	bob := ganache.Bob

// 	processor, err := processor.NewProcessor(sto, ganache.Client, ganache.Assets, ganache.Owner)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	ganache.CreateTeam("Barca", alice)
// 	ganache.CreateTeam("Madrid", bob)

// 	validUntil := big.NewInt(2000000000)
// 	playerId := big.NewInt(10)
// 	typeOfTX := uint8(1)
// 	currencyId := uint8(1)
// 	price := big.NewInt(41234)
// 	rnd := big.NewInt(42321)

// 	originOwner := ganache.GetPlayerOwner(playerId)
// 	if originOwner != ganache.Public(alice) {
// 		t.Fatalf("Expectedf originOwner ALICE but got %v", originOwner)
// 	}
// 	sto.CreateSellOrder(storage.SellOrder{
// 		PlayerId:   playerId,
// 		CurrencyId: currencyId,
// 		Price:      price,
// 		Rnd:        rnd,
// 		ValidUntil: validUntil,
// 		TypeOfTx:   typeOfTX,
// 		Signature:  "0x405c83733f474f6919032fd41bd2e37b1a3be444bc52380c0e3f4c79ce8245ce229b4b0fe3a9798b5aad5f8df5c6acc07e4810f1a111d7712bf06aee7c7384001b",
// 	})
// 	processor.Process()
// 	targetOwner := ganache.GetPlayerOwner(playerId)
// 	if targetOwner != crypto.PubkeyToAddress(alice.PublicKey) {
// 		t.Fatalf("Expectedf originOwner ALICE but got %v", targetOwner)
// 	}

// 	sto.CreateBuyOrder(storage.BuyOrder{
// 		PlayerId:  playerId,
// 		TeamId:    big.NewInt(2),
// 		Signature: "0xd36c99c9f3077a3b24d4709399f0c034bdc99e56430019ac499f46195645fbd14f742c4cee9c99a10d15177ce322850244572d48d40558728a586711ed90a3af1c",
// 	})

// 	processor.Process()
// 	targetOwner = ganache.GetPlayerOwner(playerId)
// 	if targetOwner != crypto.PubkeyToAddress(bob.PublicKey) {
// 		t.Fatalf("Expectedf originOwner BOB but got %v", targetOwner)
// 	}

// 	buyOrders, err := sto.GetBuyOrders()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(buyOrders) != 0 {
// 		t.Fatalf("Expercted 0 but got %v", len(buyOrders))
// 	}
// 	sellOrders, err := sto.GetSellOrders()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(sellOrders) != 0 {
// 		t.Fatalf("Expercted 0 but got %v", len(sellOrders))
// 	}
// }
