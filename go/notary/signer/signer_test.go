package signer_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

func TestRSV(t *testing.T) {
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

func TestRSV2(t *testing.T) {
	r, s, v, err := signer.RSV("405c83733f474f6919032fd41bd2e37b1a3be444bc52380c0e3f4c79ce8245ce229b4b0fe3a9798b5aad5f8df5c6acc07e4810f1a111d7712bf06aee7c7384001b")
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
	currencyId := uint8(1)
	price := big.NewInt(41234)
	rnd := big.NewInt(42321)
	hash, err := signer.HashPrivateMsg(
		currencyId,
		price,
		rnd,
	)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x4200de738160a9e6b8f69648fbb7feb323f73fac5acff1b7bb546bb7ac3591fa")

	sellerHiddenPrice, err := bc.Contracts.Market.HashPrivateMsg(
		&bind.CallOpts{},
		currencyId,
		price,
		rnd,
	)
	assert.NilError(t, err)
	assert.Equal(t, "0x"+hex.EncodeToString(sellerHiddenPrice[:]), hash.Hex())
}

func TestAuctionMsg(t *testing.T) {
	validUntil := uint32(2000000000)
	auctionDurationAfterOfferIsAccepted := uint32(4358487)
	playerId := big.NewInt(10)
	currencyId := uint8(1)
	price := big.NewInt(41234)
	rnd := big.NewInt(42321)

	sellDigest, err := signer.ComputeSellPlayerDigest(
		currencyId,
		price,
		rnd,
		validUntil,
		auctionDurationAfterOfferIsAccepted,
		playerId,
	)

	assert.NilError(t, err)
	result := hex.EncodeToString(sellDigest[:])
	if result != "1e5f3296caef0abef206907a7d9dea17fc035e3e5b5f23e7340c8743f7da7eae" {
		t.Fatalf("Hash error %v", result)
	}
	pvr, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	if err != nil {
		t.Fatal(err)
	}
	sig, err := signer.Sign(sellDigest.Bytes(), pvr)
	if err != nil {
		t.Fatal(err)
	}
	result = hex.EncodeToString(sig)
	if result != "c711ade1d48b670d0228eb8c6b22a1cc865ce6bce005755dacad656c976de73a14ec926dec760a0b56e2d62f9132fc81364b12d40f09c5fd23f4ff27457dd3cc1b" {
		t.Fatalf("Sign error %v", result)
	}
}

func TestPublicKeyBytesToAddress(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	assert.Equal(t, ok, true)

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := signer.PublicKeyBytesToAddress(publicKeyBytes)
	assert.Equal(t, address.Hex(), "0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7")
}

func TestHashBidMessage(t *testing.T) {
	currencyId := uint8(1)
	validUntil := uint32(2000000000)
	auctionDurationAfterOfferIsAccepted := uint32(4358487)
	playerId := big.NewInt(274877906944)
	price := big.NewInt(41234)
	auctionRnd := big.NewInt(42321)
	extraPrice := big.NewInt(0)
	bidRnd := big.NewInt(0)
	teamID := big.NewInt(274877906945)

	hash, err := signer.HashBidMessage(
		bc.Contracts.Market,
		currencyId,
		price,
		auctionRnd,
		validUntil,
		auctionDurationAfterOfferIsAccepted,
		playerId,
		extraPrice,
		bidRnd,
		teamID,
	)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(hash[:])
	if result != "aaf31b6fdf4fa69a5ebfe9874944435c0ac4c8d5c4e0f0448c56640b67789b51" {
		t.Fatalf("Hash error %v", result)
	}
	pvr, err := crypto.HexToECDSA("3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882")
	sig, err := signer.Sign(hash.Bytes(), pvr)
	if err != nil {
		t.Fatal(err)
	}
	result = hex.EncodeToString(sig)
	if result != "8c276603fcde4ec77007791c8f873eea5e63623f98c55d8277e8ece703dd38174b01cc8fe0f34bab1dd62183162e95434eb3f4ce04b096f1d912c2a0bc44beef1c" {
		t.Fatalf("Sign error %v", result)
	}
}

func TestHashBidMessage2(t *testing.T) {
	validUntil := uint32(2000000000)
	auctionDurationAfterOfferIsAccepted := uint32(4358487)
	playerId := big.NewInt(274877906944)
	currencyId := uint8(1)
	price := big.NewInt(41234)
	auctionRnd := big.NewInt(42321)
	extraPrice := big.NewInt(332)
	bidRnd := big.NewInt(1243523)
	teamID := big.NewInt(274877906945)

	auctionHash, err := signer.ComputeSellPlayerDigest(
		currencyId,
		price,
		auctionRnd,
		validUntil,
		auctionDurationAfterOfferIsAccepted,
		playerId,
	)
	assert.NilError(t, err)
	assert.Equal(t, auctionHash.Hex(), "0xa45cdd39cee0c176eac975fd5d9aae4a5185f6a53f0f4599a3f540dcf86e6c9a")

	hash, err := signer.HashBidMessageFromSellerDigest(
		bc.Contracts.Market,
		auctionHash,
		extraPrice,
		bidRnd,
		teamID,
	)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(hash[:])
	if result != "4b39c53df99608ec7fa2a43f239ff506b012323c71dec569c57136b7331c0090" {
		t.Fatalf("Hash error %v", result)
	}
	pvr, err := crypto.HexToECDSA("3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882")
	sig, err := signer.Sign(hash.Bytes(), pvr)
	if err != nil {
		t.Fatal(err)
	}
	result = hex.EncodeToString(sig)
	if result != "4f98470529ce128a0bea9a8136cbc6cc776306a4e8f375382ebd0c428825e1b718aabc4f599d58093c204390f0099f29b4d26b8495179403ec21fb2b3e023fb51c" {
		t.Fatalf("Sign error %v", result)
	}
}

func TestBidHiddenPrice(t *testing.T) {
	extraPrice := big.NewInt(332)
	buyerRandom := big.NewInt(1243523)

	hash, err := signer.BidHiddenPrice(bc.Contracts.Market, extraPrice, buyerRandom)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(hash[:])
	if result != "d46585a1479af8dcc6f2ce8495174282f8c874f1583182bbe2c9df7ae2358edc" {
		t.Fatalf("Hash error %v", result)
	}
}

func TestAuctionHiddenPrice2(t *testing.T) {
	currencyId := uint8(1)
	price := big.NewInt(41234)
	rnd := big.NewInt(4232)
	hash, err := signer.HashPrivateMsg(
		currencyId,
		price,
		rnd,
	)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x233c8f4c0d172a80f6dcca6359a08182a64d4201d33359e112e99c0025b3ed86")
}

func TestComputeSellPlayerDigest(t *testing.T) {
	currencyId := uint8(1)
	playerId := big.NewInt(11114324213423)
	price := big.NewInt(4324213423)
	rnd := big.NewInt(434324324213423)
	validUntil := uint32(235985749)
	auctionDurationAfterOfferIsAccepted := uint32(4358487)

	digest, err := signer.ComputeSellPlayerDigest(
		currencyId,
		price,
		rnd,
		validUntil,
		auctionDurationAfterOfferIsAccepted,
		playerId,
	)
	assert.NilError(t, err)
	assert.Equal(t, digest.Hex(), "0xf778fa056bd74980669505bf4666bbde172de50abe33d569f3ce597bdd81198b")

}
