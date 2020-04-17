package signer_test

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"

	"github.com/ethereum/go-ethereum/common"
	"gotest.tools/assert"
)

func TestVerifySignature(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	assert.Equal(t, ok, true)
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	data := []byte("hello")
	hash := crypto.Keccak256Hash(data)
	assert.Equal(t, hash.Hex(), "0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8")

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hexutil.Encode(signature), "0xd8ffac62ef76b5c5706f927679a86cb21df5f857f14cee09019bf270be988b5e34c29f96bb066d4a2d7a6a2a1fc8716d9c4f627b153dfbe5d8ca27b030ea211101")

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	assert.NilError(t, err)
	assert.Assert(t, bytes.Equal(sigPublicKey, publicKeyBytes))
	address := signer.PublicKeyBytesToAddress(sigPublicKey)
	assert.Equal(t, address.Hex(), "0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7")

	// sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	// assert.NilError(t, err)

	// sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	// assert.Assert(t, bytes.Equal(sigPublicKeyBytes, publicKeyBytes))

	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	assert.Assert(t, crypto.VerifySignature(sigPublicKey, hash.Bytes(), signatureNoRecoverID))
}

func TestVerifyAuctionSignature(t *testing.T) {
	validUntil := int64(2000000000)
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
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0xc50d978b8a838b6c437a162a94c715f95e92e11fe680cf0f1caf054ad78cd796")

	privateKey, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "075ddf60b307abf0ecf323dcdd57230fcb81b30217fb947ee5dbd683cb8bcf074a63f87c97c736f85cd3e56e95f4fcc1e9b159059817915d0be68f944f5b4e5301")

	valid, err := signer.VerifySignature(hash.Bytes(), signature)
	assert.NilError(t, err)
	assert.Assert(t, valid)

	address, err := signer.AddressFromSignature(hash.Bytes(), signature)
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7")
}

func TestSignatureInvalidId(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)

	hash := common.HexToHash("0xc50d978b8a838b6c437a162a94c715f95e92e11fe680cf0f1caf054ad78cd796")
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "075ddf60b307abf0ecf323dcdd57230fcb81b30217fb947ee5dbd683cb8bcf074a63f87c97c736f85cd3e56e95f4fcc1e9b159059817915d0be68f944f5b4e531c")

	_, err = signer.VerifySignature(hash.Bytes(), signature)
	assert.Error(t, err, "invalid signature recovery id")
}
