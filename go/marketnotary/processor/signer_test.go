package processor_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/processor"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestRSV(t *testing.T) {
	signer := processor.NewSigner(nil)
	_, _, _, err := signer.RSV("0x0")
	if err == nil {
		t.Fatal("No error on wrong signature")
	}
	r, s, v, err := signer.RSV("0x405c83733f474f6919032fd41bd2e37b1a3be444bc52380c0e3f4c79ce8245ce229b4b0fe3a9798b5aad5f8df5c6acc07e4810f1a111d7712bf06aee7c7384001b")
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

func TestAuctionHiddenPrice(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}

	err = bc.DeployContracts(bc.Owner)
	signer := processor.NewSigner(bc.Market)
	currencyId := uint8(1)
	price := big.NewInt(41234)
	rnd := big.NewInt(42321)
	privHash, err := signer.HashPrivateMsg(
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
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}

	err = bc.DeployContracts(bc.Owner)
	signer := processor.NewSigner(bc.Market)
	validUntil := big.NewInt(2000000000)
	playerId := big.NewInt(10)
	currencyId := uint8(1)
	price := big.NewInt(41234)
	rnd := big.NewInt(42321)

	hash, err := signer.HashSellMessage(
		currencyId,
		price,
		rnd,
		validUntil,
		playerId,
	)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(hash[:])
	if result != "c50d978b8a838b6c437a162a94c715f95e92e11fe680cf0f1caf054ad78cd796" {
		t.Fatalf("Hash error %v", result)
	}
}

// func TestHashAgreeToBuyMessage(t *testing.T) {
// 	ganache := testutils.NewGanache()
// 	signer := processor.NewSigner(ganache.Market)
// 	validUntil := big.NewInt(2000000000)
// 	playerId := big.NewInt(10)
// 	typeOfTx := uint8(1)
// 	currencyId := uint8(1)
// 	price := big.NewInt(41234)
// 	rnd := big.NewInt(42321)

// 	hash, err := signer.HashBuyMessage(
// 		currencyId,
// 		price,
// 		rnd,
// 		validUntil,
// 		playerId,
// 		typeOfTx,
// 		big.NewInt(2),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result := hex.EncodeToString(hash[:])
// 	if result != "0d84fd72fb639204abba9869b3fcb7855df4b83c121c1d6fd679f90c828d5528" {
// 		t.Fatalf("Hash error %v", result)
// 	}
// }
