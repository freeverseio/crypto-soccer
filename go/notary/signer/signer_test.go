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
	hash, err := signer.HidePrice(
		currencyId,
		price,
		rnd,
	)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x4200de738160a9e6b8f69648fbb7feb323f73fac5acff1b7bb546bb7ac3591fa")

	bcHash, err := bc.Contracts.Market.HashPrivateMsg(
		&bind.CallOpts{},
		currencyId,
		price,
		rnd,
	)
	assert.NilError(t, err)
	assert.Equal(t, "0x"+hex.EncodeToString(bcHash[:]), hash.Hex())
}

// the hashes in this test comes from from JS-Solidity test ref XX1
func TestAuctionMsg(t *testing.T) {
	validUntil := int64(235985749)
	offerValidUntil := int64(4358487)
	playerId := big.NewInt(11114324213423)
	currencyId := uint8(1)
	price := big.NewInt(345)
	rnd := big.NewInt(1234)

	hash, err := signer.ComputePutAssetForSaleDigest(
		currencyId,
		price,
		rnd,
		validUntil,
		offerValidUntil,
		playerId,
	)
	assert.NilError(t, err)
	result := hex.EncodeToString(hash[:])
	if result != "376b87a3db2c3ef6e1189a96303454a32fd8bf21bfe0a470e68be98e57d36495" {
		t.Fatalf("Hash error %v", result)
	}
	pvr, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	if err != nil {
		t.Fatal(err)
	}
	sig, err := signer.Sign(hash.Bytes(), pvr)
	if err != nil {
		t.Fatal(err)
	}
	result = hex.EncodeToString(sig)
	if result != "f0e4f8fe6502bb950fa45283832d117dda9876e1bf92c29808ab9072fd717cc3756ee55cd659cc33ed2d3d0aa6f290f3f583045e9b91c32cab64747b8b43c7701b" {
		t.Fatalf("Sign error %v", result)
	}
	// auction Id 1:
	auctionId, err2 := signer.ComputeAuctionId(
		currencyId,
		price,
		rnd,
		validUntil,
		0,
		playerId,
	)
	assert.NilError(t, err2)
	result = hex.EncodeToString(auctionId[:])
	if result != "03214d89eb62587cbb48c9056dba878f839a4ebad3ad75f8826d76c566e4acd0" {
		t.Fatalf("Hash error %v", result)
	}
	// auction Id 2:
	auctionId, err2 = signer.ComputeAuctionId(
		currencyId,
		price,
		rnd,
		validUntil,
		offerValidUntil,
		playerId,
	)
	assert.NilError(t, err2)
	result = hex.EncodeToString(auctionId[:])
	if result != "f06dfe068a4aa5621dddc8d424ca97c0bd6a2ef5e9af94ba6ba3550beb6e0438" {
		t.Fatalf("Hash error %v", result)
	}
	// auction Id 3:
	auctionId, err2 = signer.ComputeAuctionId(
		currencyId,
		price,
		rnd,
		0,
		offerValidUntil,
		playerId,
	)
	assert.NilError(t, err2)
	result = hex.EncodeToString(auctionId[:])
	if result != "f06dfe068a4aa5621dddc8d424ca97c0bd6a2ef5e9af94ba6ba3550beb6e0438" {
		t.Fatalf("Hash error %v", result)
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

// the hashes from this test come from Solidity+JS, ref XX2
func TestHashBidMessage(t *testing.T) {
	validUntil := int64(2000000000)
	offerValidUntil := validUntil - 300
	playerId := big.NewInt(274877906944)
	currencyId := uint8(1)
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
		offerValidUntil,
		playerId,
		extraPrice,
		bidRnd,
		teamID,
	)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(hash[:])
	if result != "c272356ecde3af2f8c3d04879f8c32b1ccc091ff155f0c9d2412f35620ccd235" {
		t.Fatalf("Hash error %v", result)
	}
	pvr, err := crypto.HexToECDSA("3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882")
	sig, err := signer.Sign(hash.Bytes(), pvr)
	if err != nil {
		t.Fatal(err)
	}
	result = hex.EncodeToString(sig)
	if result != "b3a17e9c821e4e6bd12554168add461d77f15da42657d284eb9a459ee646bdd15024a60a6332f8617ff82f5a1c5fedf4e005da620ac01673c71ad5c08bba2d221c" {
		t.Fatalf("Sign error %v", result)
	}
}

func TestHasBidFromAuctionId(t *testing.T) {
	validUntil := int64(2000000000)
	offerValidUntil := validUntil - 300
	playerId := big.NewInt(274877906944)
	currencyId := uint8(1)
	price := big.NewInt(41234)
	auctionRnd := big.NewInt(42321)
	extraPrice := big.NewInt(0)
	bidRnd := big.NewInt(0)
	teamID := big.NewInt(274877906945)

	auctionId, err := signer.ComputeAuctionId(
		currencyId,
		price,
		auctionRnd,
		validUntil,
		offerValidUntil,
		playerId,
	)
	assert.NilError(t, err)

	hash, err := signer.HasBidFromAuctionId(
		bc.Contracts.Market,
		auctionId,
		extraPrice,
		bidRnd,
		teamID,
	)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(hash[:])
	if result != "c272356ecde3af2f8c3d04879f8c32b1ccc091ff155f0c9d2412f35620ccd235" {
		t.Fatalf("Hash error %v", result)
	}
	pvr, err := crypto.HexToECDSA("3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882")
	sig, err := signer.Sign(hash.Bytes(), pvr)
	if err != nil {
		t.Fatal(err)
	}
	result = hex.EncodeToString(sig)
	if result != "b3a17e9c821e4e6bd12554168add461d77f15da42657d284eb9a459ee646bdd15024a60a6332f8617ff82f5a1c5fedf4e005da620ac01673c71ad5c08bba2d221c" {
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
	hash, err := signer.HidePrice(
		currencyId,
		price,
		rnd,
	)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x233c8f4c0d172a80f6dcca6359a08182a64d4201d33359e112e99c0025b3ed86")
}
