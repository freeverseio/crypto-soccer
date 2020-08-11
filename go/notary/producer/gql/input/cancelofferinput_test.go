package input_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

func TestCancelOfferHash(t *testing.T) {
	in := input.CancelOfferInput{}
	in.OfferId = ""
	hash, err := in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x4b60b7c7b3f67bb245d5360199fe2754fff8a649a3b483d945f0a77e9897072b")
	in.OfferId = "43"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x5757964f4d77a1fdc41d891587d2dd6fd593df7d2933a5d8f2ecab7ddf26c6fe")
}
func TestCancelOfferValidSignature(t *testing.T) {
	in := input.CancelOfferInput{}
	in.OfferId = "4345345341"
	in.Signature = "075ddf60b307abf0ecf323dcdd57230fcb81b30217fb947ee5dbd683cb8bcf074a63f87c97c736f85cd3e56e95f4fcc1e9b159059817915d0be68f944f5b4e531c"
	valid, err := in.VerifySignature()
	assert.NilError(t, err)
	assert.Assert(t, valid)
}

func TestCancelOfferGetSigner(t *testing.T) {
	in := input.CancelOfferInput{}
	in.OfferId = "c50d978b8a838b6c437a162a94c715f95e92e11fe680cf0f1caf054ad78cd796"

	hash, err := in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x20d4c8848f2c767dbe7dc79e56e05e61d717a27ec94d635d2ef888f20ed7335c")

	pvc, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	sign, err := signer.Sign(hash.Bytes(), pvc)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "ae2431f4d5e8d8f05b3478bbaa293213c697c3d3ef09ff02b3a9b2ffb98199b25622dc55c1774809276149caac35cd1ccef358578a7c7b2aabd7ec0a15b017b81c")

	address, err := in.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
}
