package input_test

import (
	"encoding/hex"
	"strconv"
	"testing"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

func TestAcceptOfferInputDigestAndID(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = "2000000000"
	in.AuctionDurationAfterOfferIsAccepted = "3600"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	digest, err := in.Digest()
	assert.NilError(t, err)
	expected := "852d3f820972df98e3cb66bc87c118b36d302c9906fa11eed4ff54cb2235f362"
	assert.Equal(t, digest.Hex(), "0x"+expected)
	id, err := in.AuctionID()
	assert.NilError(t, err)
	assert.Equal(t, string(id), expected)
}
func TestAcceptOfferValidSignature(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = "2000000000"
	in.AuctionDurationAfterOfferIsAccepted = "3600"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321

	in.Signature = "075ddf60b307abf0ecf323dcdd57230fcb81b30217fb947ee5dbd683cb8bcf074a63f87c97c736f85cd3e56e95f4fcc1e9b159059817915d0be68f944f5b4e531c"
	valid, err := in.VerifySignature()
	assert.NilError(t, err)
	assert.Assert(t, valid)
}

func TestAcceptOfferSignerAddress(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = "235985749"
	in.AuctionDurationAfterOfferIsAccepted = "4358487"
	in.PlayerId = "11114324213423"
	in.CurrencyId = 1
	in.Price = 345
	in.Rnd = 1234
	in.Signature = "f0e4f8fe6502bb950fa45283832d117dda9876e1bf92c29808ab9072fd717cc3756ee55cd659cc33ed2d3d0aa6f290f3f583045e9b91c32cab64747b8b43c7701b"
	address, err := in.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7")
}

func TestAcceptOfferIsSignerOwner(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = "2000000000"
	in.AuctionDurationAfterOfferIsAccepted = "3600"
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321

	digest, err := in.Digest()
	assert.NilError(t, err)
	signature, err := signer.Sign(digest.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	isOwner, err := in.IsSignerOwner(*bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, isOwner)
}

func TestAcceptOfferIsValidBlockchain(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	in.AuctionDurationAfterOfferIsAccepted = "3600"
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321

	digest, err := in.Digest()
	assert.NilError(t, err)
	signature, err := signer.Sign(digest.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	isValid, err := in.IsValidForBlockchain(*bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, isValid)
}
