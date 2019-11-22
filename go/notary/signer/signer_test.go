package signer_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestRSV(t *testing.T) {
	signer := signer.NewSigner(nil, nil)
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
	signer := signer.NewSigner(bc.Market, nil)
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

func TestAuctionMsg(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}

	err = bc.DeployContracts(bc.Owner)
	signer := signer.NewSigner(bc.Market, nil)
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
	pvr, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	if err != nil {
		t.Fatal(err)
	}
	sig, err := signer.Sign(hash, pvr)
	if err != nil {
		t.Fatal(err)
	}
	result = hex.EncodeToString(sig)
	if result != "075ddf60b307abf0ecf323dcdd57230fcb81b30217fb947ee5dbd683cb8bcf074a63f87c97c736f85cd3e56e95f4fcc1e9b159059817915d0be68f944f5b4e531c" {
		t.Fatalf("Sign error %v", result)
	}
}

func TestHashBidMessage(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	signer := signer.NewSigner(bc.Market, nil)

	validUntil := big.NewInt(2000000000)
	playerId := big.NewInt(274877906944)
	currencyId := uint8(1)
	price := big.NewInt(41234)
	auctionRnd := big.NewInt(42321)
	extraPrice := big.NewInt(332)
	bidRnd := big.NewInt(1243523)
	teamID := big.NewInt(274877906945)
	isOffer2StartAuction := true

	hash, err := signer.HashBidMessage(
		currencyId,
		price,
		auctionRnd,
		validUntil,
		playerId,
		extraPrice,
		bidRnd,
		teamID,
		isOffer2StartAuction,
	)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(hash[:])
	if result != "e04d23ec0424b22adec87879118715ce75997a4fd47897c398f3a8cad79b3041" {
		t.Fatalf("Hash error %v", result)
	}
	pvr, err := crypto.HexToECDSA("3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882")
	sig, err := signer.Sign(hash, pvr)
	if err != nil {
		t.Fatal(err)
	}
	result = hex.EncodeToString(sig)
	if result != "dbe104e7b51c9b1e38cdda4e31c2036e531f7d3338d392bee2f526c4c892437f5e50ddd44224af8b3bd92916b93e4b0d7af2974175010323da7dedea19f30d621c" {
		t.Fatalf("Sign error %v", result)
	}
}

func TestBidHiddenPrice(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	signer := signer.NewSigner(bc.Market, nil)
	extraPrice := big.NewInt(332)
	buyerRandom := big.NewInt(1243523)

	hash, err := signer.BidHiddenPrice(extraPrice, buyerRandom)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(hash[:])
	if result != "d46585a1479af8dcc6f2ce8495174282f8c874f1583182bbe2c9df7ae2358edc" {
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
