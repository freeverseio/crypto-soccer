package input_test

import (
	"encoding/hex"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

func TestCancelAuctionHash(t *testing.T) {
	in := input.CancelAuctionInput{}
	in.AuctionId = ""
	hash, err := in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x4b60b7c7b3f67bb245d5360199fe2754fff8a649a3b483d945f0a77e9897072b")
	in.AuctionId = "43"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x5a24d5cbf413a599aac8109527f6af95eedc33b69f57f23e384ef20df5e9651e")
}
func TestCancelAuctionValidSignature(t *testing.T) {
	in := input.CancelAuctionInput{}
	in.AuctionId = "434534534"
	in.Signature = "075ddf60b307abf0ecf323dcdd57230fcb81b30217fb947ee5dbd683cb8bcf074a63f87c97c736f85cd3e56e95f4fcc1e9b159059817915d0be68f944f5b4e531c"
	valid, err := in.VerifySignature()
	assert.NilError(t, err)
	assert.Assert(t, valid)
}

func TestCancelAuctionGetSigner(t *testing.T) {
	in := input.CancelAuctionInput{}
	in.AuctionId = "c50d978b8a838b6c437a162a94c715f95e92e11fe680cf0f1caf054ad78cd796"

	hash, err := in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0xdebecc786086a1cf7be2b6d5d7f90bf93184f13239004fe5fe4dc2c30bc825ea")

	sign, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "63288c53fa0a7be0f33b0bc39c3003ec133b3aec48cfe00904101cc8612750f15ed9e786430a06b7d5c5803b1c8c5220ed20a286817dabc1a37c5b4c18832e011b")

	address, err := in.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
}
