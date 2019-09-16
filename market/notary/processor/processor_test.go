package processor_test

import (
	"encoding/hex"
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

	timezone := uint8(1)
	countryIdxInTZ := big.NewInt(0)
	teamId0, err := ganache.Market.EncodeTZCountryAndVal(&bind.CallOpts{}, timezone, countryIdxInTZ, big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ganache.Market.TransferBotToAddr(bind.NewKeyedTransactor(alice), teamId0, crypto.PubkeyToAddress(alice.PublicKey))
	if err != nil {
		t.Fatal(err)
	}
	teamId1, err := ganache.Market.EncodeTZCountryAndVal(&bind.CallOpts{}, timezone, countryIdxInTZ, big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ganache.Market.TransferBotToAddr(bind.NewKeyedTransactor(bob), teamId1, crypto.PubkeyToAddress(bob.PublicKey))
	if err != nil {
		t.Fatal(err)
	}
	teamId2, err := ganache.Market.EncodeTZCountryAndVal(&bind.CallOpts{}, timezone, countryIdxInTZ, big.NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ganache.Market.TransferBotToAddr(bind.NewKeyedTransactor(alice), teamId2, crypto.PubkeyToAddress(alice.PublicKey))
	if err != nil {
		t.Fatal(err)
	}
	team0PlayerIds, err := ganache.Market.GetPlayerIdsInTeam(&bind.CallOpts{}, teamId0)
	if err != nil {
		t.Fatal(err)
	}
	playerId := team0PlayerIds[0]
	owner := ganache.GetPlayerOwner(playerId)
	if owner != ganache.Public(alice) {
		t.Fatalf("Expected owner ALICE but got %v", owner)
	}
	err = ganache.TransferPlayer(playerId, teamId1)
	if err != nil {
		t.Fatal(err)
	}
	owner = ganache.GetPlayerOwner(playerId)
	if owner != ganache.Public(bob) {
		t.Fatalf("Expected owner BOB but got %v", owner.Hex())
	}
	err = ganache.TransferPlayer(playerId, teamId2)
	if err != nil {
		t.Fatal(err)
	}
	owner = ganache.GetPlayerOwner(playerId)
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

	timezone := uint8(1)
	countryIdxInTZ := big.NewInt(0)
	teamId0, err := ganache.Market.EncodeTZCountryAndVal(&bind.CallOpts{}, timezone, countryIdxInTZ, big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ganache.Market.TransferBotToAddr(bind.NewKeyedTransactor(alice), teamId0, crypto.PubkeyToAddress(alice.PublicKey))
	if err != nil {
		t.Fatal(err)
	}
	teamId1, err := ganache.Market.EncodeTZCountryAndVal(&bind.CallOpts{}, timezone, countryIdxInTZ, big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ganache.Market.TransferBotToAddr(bind.NewKeyedTransactor(bob), teamId1, crypto.PubkeyToAddress(bob.PublicKey))
	if err != nil {
		t.Fatal(err)
	}
	team0PlayerIds, err := ganache.Market.GetPlayerIdsInTeam(&bind.CallOpts{}, teamId0)
	if err != nil {
		t.Fatal(err)
	}

	validUntil := big.NewInt(2000000000)
	playerId := team0PlayerIds[0]
	typeOfTX := uint8(1)
	currencyId := uint8(1)
	price := big.NewInt(41234)
	rnd := big.NewInt(42321)
	teamId := teamId1

	originOwner := ganache.GetPlayerOwner(playerId)
	if originOwner != ganache.Public(alice) {
		t.Fatalf("Expected originOwner ALICE but got %v", originOwner)
	}
	sto.CreateSellOrder(storage.SellOrder{
		PlayerId:   playerId,
		CurrencyId: currencyId,
		Price:      price,
		Rnd:        rnd,
		ValidUntil: validUntil,
		TypeOfTx:   typeOfTX,
		Signature:  "0xac466c2139f6edce74d18161252922d8368dce25c3e508de98e8659e9a994a000dd33bd3034aea26fe99b54b1df240041f77afb0a2be508a83e7d35482b20a951c",
	})
	processor.Process()
	targetOwner := ganache.GetPlayerOwner(playerId)
	if targetOwner != crypto.PubkeyToAddress(alice.PublicKey) {
		t.Fatalf("Expectedf originOwner ALICE but got %v", targetOwner)
	}

	sto.CreateBuyOrder(storage.BuyOrder{
		PlayerId:  playerId,
		TeamId:    teamId,
		Signature: "0x44bb117064e1e2a8ef5fed99f3ec9281f95ef7caea595db2c36071963f74e4c904e8c61d6cb75aaef61718e1d2dff49bc3c55c886e7b3d3e73db31a1af3c61721b",
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
