package processor_test

import (
	"encoding/hex"
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

	ganache.AssignTeam(big.NewInt(0), alice)
	ganache.AssignTeam(big.NewInt(1), bob)

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
	ganache.AssignTeam(big.NewInt(3), alice)
	err = ganache.TransferPlayer(player, big.NewInt(3))
	if err != nil {
		t.Fatal(err)
	}
	owner = ganache.GetPlayerOwner(player)
	if owner != ganache.Public(alice) {
		t.Fatalf("Expected owner ALICE but got %v", owner)
	}
}

func TestRSV(t *testing.T) {
	r, s, v, err := processor.RSV("0x405c83733f474f6919032fd41bd2e37b1a3be444bc52380c0e3f4c79ce8245ce229b4b0fe3a9798b5aad5f8df5c6acc07e4810f1a111d7712bf06aee7c7384001b")
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(r[:])
	if result != "405c83733f474f6919032fd41bd2e37b1a3be444bc52380c0e3f4c79ce8245ce" {
		t.Fatalf("r error %v", result)
	}
	result = hex.EncodeToString(s[:])
	if result != "229b4b0fe3a9798b5aad5f8df5c6acc07e4810f1a111d7712bf06aee7c738400" {
		t.Fatalf("s error %v", result)
	}
	if v != 0x1b {
		t.Fatalf("Error in v %v", v)
	}
}

func TestHashPRivateMessage(t *testing.T) {
	ganache := testutils.NewGanache()
	processor, err := processor.NewProcessor(nil, ganache.Client, ganache.Market, ganache.Owner)
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

func TestBuildPutForSaleMessage(t *testing.T) {
	ganache := testutils.NewGanache()
	processor, err := processor.NewProcessor(nil, ganache.Client, ganache.Market, ganache.Owner)
	if err != nil {
		t.Fatal(err)
	}
	validUntil := big.NewInt(2000000000)
	playerId := big.NewInt(10)
	typeOfTx := uint8(1)
	currencyId := uint8(1)
	price := big.NewInt(41234)
	rnd := big.NewInt(42321)

	hash, err := processor.HashSellMessage(
		currencyId,
		price,
		rnd,
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

func TestHashAgreeToBuyMessage(t *testing.T) {
	ganache := testutils.NewGanache()
	processor, err := processor.NewProcessor(nil, ganache.Client, ganache.Market, ganache.Owner)
	if err != nil {
		t.Fatal(err)
	}
	validUntil := big.NewInt(2000000000)
	playerId := big.NewInt(10)
	typeOfTx := uint8(1)
	currencyId := uint8(1)
	price := big.NewInt(41234)
	rnd := big.NewInt(42321)

	hash, err := processor.HashBuyMessage(
		currencyId,
		price,
		rnd,
		validUntil,
		playerId,
		typeOfTx,
		big.NewInt(2),
	)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(hash[:])
	if result != "0d84fd72fb639204abba9869b3fcb7855df4b83c121c1d6fd679f90c828d5528" {
		t.Fatalf("Hash error %v", result)
	}
}

func TestProcess(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache := testutils.NewGanache()
	alice, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	bob, err := crypto.HexToECDSA("3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882")

	if err != nil {
		t.Fatal(err)
	}
	value := new(big.Int)
	value.SetString("50000000000000000000", 10)
	_, err = ganache.TransferWei(value, ganache.Owner, ganache.Public(alice))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ganache.TransferWei(value, ganache.Owner, ganache.Public(bob))
	if err != nil {
		t.Fatal(err)
	}

	processor, err := processor.NewProcessor(sto, ganache.Client, ganache.Market, ganache.Owner)
	if err != nil {
		t.Fatal(err)
	}

	ganache.AssignTeam(big.NewInt(0), alice)
	ganache.AssignTeam(big.NewInt(1), bob)

	validUntil := big.NewInt(2000000000)
	playerId := big.NewInt(10)
	typeOfTX := uint8(1)
	currencyId := uint8(1)
	price := big.NewInt(41234)
	rnd := big.NewInt(42321)
	teamId := big.NewInt(2)

	originOwner := ganache.GetPlayerOwner(playerId)
	if originOwner != ganache.Public(alice) {
		t.Fatalf("Expectedf originOwner ALICE but got %v", originOwner)
	}
	sto.CreateSellOrder(storage.SellOrder{
		PlayerId:   playerId,
		CurrencyId: currencyId,
		Price:      price,
		Rnd:        rnd,
		ValidUntil: validUntil,
		TypeOfTx:   typeOfTX,
		Signature:  "0x405c83733f474f6919032fd41bd2e37b1a3be444bc52380c0e3f4c79ce8245ce229b4b0fe3a9798b5aad5f8df5c6acc07e4810f1a111d7712bf06aee7c7384001b",
	})
	processor.Process()
	targetOwner := ganache.GetPlayerOwner(playerId)
	if targetOwner != crypto.PubkeyToAddress(alice.PublicKey) {
		t.Fatalf("Expectedf originOwner ALICE but got %v", targetOwner)
	}

	sto.CreateBuyOrder(storage.BuyOrder{
		PlayerId:  playerId,
		TeamId:    teamId,
		Signature: "0x0f998640c4c2348dfcd0077be8673b34ce716e02af35d65792614294759a9bc26951b536f0dc481a9e1e4642bcf18a692d1bf673d911589031106f634df42cca1b",
	})

	processor.Process()
	targetOwner = ganache.GetPlayerOwner(playerId)
	if targetOwner != crypto.PubkeyToAddress(bob.PublicKey) {
		t.Fatalf("Expected originOwner BOB but got %v", targetOwner)
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

func TestProcess2(t *testing.T) {
	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	ganache := testutils.NewGanache()
	alice, err := crypto.HexToECDSA("431c142d0c66ffae950717d298066b338a8e43fd598a9e24211ca3c64606349e")
	bob, err := crypto.HexToECDSA("3c53fcc079dd1cbc921d7f2a739ccbc47ad35a54f79c072eb9a002c764766e71")

	if err != nil {
		t.Fatal(err)
	}
	value := new(big.Int)
	value.SetString("50000000000000000000", 10)
	_, err = ganache.TransferWei(value, ganache.Owner, ganache.Public(alice))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ganache.TransferWei(value, ganache.Owner, ganache.Public(bob))
	if err != nil {
		t.Fatal(err)
	}

	proc, err := processor.NewProcessor(sto, ganache.Client, ganache.Market, ganache.Owner)
	if err != nil {
		t.Fatal(err)
	}

	ganache.AssignTeam(big.NewInt(0), alice)
	ganache.AssignTeam(big.NewInt(1), bob)
	ganache.AssignTeam(big.NewInt(2), bob)

	validUntil := big.NewInt(156836459600)
	playerId := big.NewInt(1)
	typeOfTX := uint8(1)
	currencyId := uint8(1)
	price := big.NewInt(4)
	rnd := big.NewInt(1988006456)
	teamId := big.NewInt(3)
	const singatureSeller = "0xbc4a5732af32c022c68ff8ca8d314ef49ec43b415b04233471cdbfc81e979eb7428fe7c411e9c7c315ff733081794925e781a54810006e7f7baf3683144613821b"
	const singatureBuyer = "0x4ba63c8cb6315fd75658eb193a2f85c6d5114b5436caef42ecfa7188909ed6297a63c8178d964b2b16c5599c885020fe2ec04870f6ee3ed6b4b2da001d961c8d1c"

	originOwner := ganache.GetPlayerOwner(playerId)
	if originOwner != ganache.Public(alice) {
		t.Fatalf("Expectedf originOwner ALICE but got %v", originOwner)
	}
	sto.CreateSellOrder(storage.SellOrder{
		PlayerId:   playerId,
		CurrencyId: currencyId,
		Price:      price,
		Rnd:        rnd,
		ValidUntil: validUntil,
		TypeOfTx:   typeOfTX,
		Signature:  singatureSeller,
	})
	proc.Process()
	targetOwner := ganache.GetPlayerOwner(playerId)
	if targetOwner != crypto.PubkeyToAddress(alice.PublicKey) {
		t.Fatalf("Expectedf originOwner ALICE but got %v", targetOwner)
	}

	sto.CreateBuyOrder(storage.BuyOrder{
		PlayerId:  playerId,
		TeamId:    teamId,
		Signature: singatureBuyer,
	})

	proc.Process()
	targetOwner = ganache.GetPlayerOwner(playerId)
	if targetOwner != crypto.PubkeyToAddress(bob.PublicKey) {
		t.Fatalf("Expected originOwner BOB but got %v", targetOwner)
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

// func TestProcess3(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// ganache := testutils.NewGanache()
// 	// alice, err := crypto.HexToECDSA("431c142d0c66ffae950717d298066b338a8e43fd598a9e24211ca3c64606349e")
// 	// bob, err := crypto.HexToECDSA("3c53fcc079dd1cbc921d7f2a739ccbc47ad35a54f79c072eb9a002c764766e71")

// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	// value := new(big.Int)
// 	// value.SetString("50000000000000000000", 10)
// 	// _, err = ganache.TransferWei(value, ganache.Owner, ganache.Public(alice))
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	// _, err = ganache.TransferWei(value, ganache.Owner, ganache.Public(bob))
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }

// 	// log.Info("Dial the Ethereum client: ", ethereumClient)
// 	client, err := ethclient.Dial("https://devnet.busyverse.com/web3")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// t.Info("Creating Assets bindings to: ", assetsContractAddress)
// 	assetsContract, err := assets.NewAssets(common.HexToAddress("0xE5094517AeE4f34811838ef7493abe0527e3B2F5"), client)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	privateKey, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	processor, err := processor.NewProcessor(sto, client, assetsContract, privateKey)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// ganache.CreateTeam("Venezie", alice)
// 	// ganache.CreateTeam("Barca", bob)
// 	// ganache.CreateTeam("Madrid", bob)

// 	validUntil := big.NewInt(156836459600)
// 	playerId := big.NewInt(1)
// 	typeOfTX := uint8(1)
// 	currencyId := uint8(1)
// 	price := big.NewInt(4)
// 	rnd := big.NewInt(1988006456)
// 	teamId := big.NewInt(3)

// 	// originOwner := ganache.GetPlayerOwner(playerId)
// 	// if originOwner != ganache.Public(alice) {
// 	// 	t.Fatalf("Expectedf originOwner ALICE but got %v", originOwner)
// 	// }
// 	sto.CreateSellOrder(storage.SellOrder{
// 		PlayerId:   playerId,
// 		CurrencyId: currencyId,
// 		Price:      price,
// 		Rnd:        rnd,
// 		ValidUntil: validUntil,
// 		TypeOfTx:   typeOfTX,
// 		Signature:  "0xbc4a5732af32c022c68ff8ca8d314ef49ec43b415b04233471cdbfc81e979eb7428fe7c411e9c7c315ff733081794925e781a54810006e7f7baf3683144613821b",
// 	})
// 	processor.Process()
// 	// targetOwner := ganache.GetPlayerOwner(playerId)
// 	// if targetOwner != crypto.PubkeyToAddress(alice.PublicKey) {
// 	// 	t.Fatalf("Expectedf originOwner ALICE but got %v", targetOwner)
// 	// }

// 	sto.CreateBuyOrder(storage.BuyOrder{
// 		PlayerId:  playerId,
// 		TeamId:    teamId,
// 		Signature: "0x4ba63c8cb6315fd75658eb193a2f85c6d5114b5436caef42ecfa7188909ed6297a63c8178d964b2b16c5599c885020fe2ec04870f6ee3ed6b4b2da001d961c8d1c",
// 	})

// 	err = processor.Process()
// 	if err != nil {
// 		t.Fatal(err)

// 	}
// 	t.Fatal("ciao")
// 	// targetOwner = ganache.GetPlayerOwner(playerId)
// 	// if targetOwner != crypto.PubkeyToAddress(bob.PublicKey) {
// 	// 	t.Fatalf("Expected originOwner BOB but got %v", targetOwner)
// 	// }

// 	// buyOrders, err := sto.GetBuyOrders()
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	// if len(buyOrders) != 0 {
// 	// 	t.Fatalf("Expercted 0 but got %v", len(buyOrders))
// 	// }
// 	// sellOrders, err := sto.GetSellOrders()
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	// if len(sellOrders) != 0 {
// 	// 	t.Fatalf("Expercted 0 but got %v", len(sellOrders))
// 	// }
// }
