package signer_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"

	"gotest.tools/assert"
)

func TestVerifyAuctionSignature(t *testing.T) {
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
	assert.Equal(t, hash.Hex(), "0x376b87a3db2c3ef6e1189a96303454a32fd8bf21bfe0a470e68be98e57d36495")
	assert.Equal(t, hex.EncodeToString(hash[:]), "376b87a3db2c3ef6e1189a96303454a32fd8bf21bfe0a470e68be98e57d36495")

	alice, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)

	// note that cryptp.Sign and signer.Sign differ in the last hex value! compare to signer_test, ref XX1
	signature, err := crypto.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "f0e4f8fe6502bb950fa45283832d117dda9876e1bf92c29808ab9072fd717cc3756ee55cd659cc33ed2d3d0aa6f290f3f583045e9b91c32cab64747b8b43c77000")

	signerFromSigAndHash, err := signer.AddressFromHashAndSignature(hash.Bytes(), signature)
	assert.NilError(t, err)
	assert.Equal(t, signerFromSigAndHash.Hex(), crypto.PubkeyToAddress(alice.PublicKey).Hex())

}
