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

func TestAcceptOfferInputHash(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = "2000000000"
	in.OfferValidUntil = "199999000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	hash, err := in.SellerDigest()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0xe30150bc666d8a20396c30794d1d1eaf86ecf427d4f4e3d8a4aa87d4aa3fc4b5")
}
func TestAcceptOfferInputAuctionID(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = "2000000000"
	in.OfferValidUntil = "199999000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	id, err := in.AuctionID()
	assert.NilError(t, err)
	assert.Equal(t, string(id), "0xe30150bc666d8a20396c30794d1d1eaf86ecf427d4f4e3d8a4aa87d4aa3fc4b5")
}
func TestAcceptOfferValidSignature(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = "2000000000"
	in.OfferValidUntil = "199999000"
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
	in.ValidUntil = "2000000000"
	in.OfferValidUntil = "199999000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	in.Signature = "075ddf60b307abf0ecf323dcdd57230fcb81b30217fb947ee5dbd683cb8bcf074a63f87c97c736f85cd3e56e95f4fcc1e9b159059817915d0be68f944f5b4e531c"
	address, err := in.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7")
}

func TestAcceptOfferIsSignerOwner(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = "2000000000"
	in.OfferValidUntil = "199999000"
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.SellerDigest()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	isOwner, err := in.IsSignerOwner(*bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, isOwner)
}

func TestAcceptOfferIsValidBlockchain(t *testing.T) {
	in := input.AcceptOfferInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+2000, 10)
	in.OfferValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.SellerDigest()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	isValid, err := in.IsValidForBlockchain(*bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, isValid)
}
