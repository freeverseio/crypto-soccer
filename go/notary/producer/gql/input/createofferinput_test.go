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

func TestCreateOfferInputHash(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	in.TeamId = "20"
	hash, err := in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x1cb2cbf8063631f3c7968aa0efd48df796f7f512faea66effd5c48f36a0660ce")
}

func TestCreateOfferInputID(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	in.TeamId = "20"
	id, err := in.ID()
	assert.NilError(t, err)
	assert.Equal(t, string(id), "1cb2cbf8063631f3c7968aa0efd48df796f7f512faea66effd5c48f36a0660ce")
}

func TestCreateOfferValidSignature(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.TeamId = "20"
	in.Rnd = 42321

	in.Signature = "075ddf60b307abf0ecf323dcdd57230fcb81b30217fb947ee5dbd683cb8bcf074a63f87c97c736f85cd3e56e95f4fcc1e9b159059817915d0be68f944f5b4e531c"
	valid, err := in.VerifySignature()
	assert.NilError(t, err)
	assert.Assert(t, valid)
}

func TestCreateOfferSignerAddress(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.TeamId = "20"
	in.Rnd = 42321
	in.Signature = "075ddf60b307abf0ecf323dcdd57230fcb81b30217fb947ee5dbd683cb8bcf074a63f87c97c736f85cd3e56e95f4fcc1e9b159059817915d0be68f944f5b4e531c"
	address, err := in.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0x15b86743899dDE665a5085E257283824ea11a27B")
}

func TestCreateOfferIsSignerOwner(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "274877906940"
	in.CurrencyId = 1
	in.TeamId = "20"
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	isOwner, err := in.IsSignerOwner(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, isOwner, false)
}

func TestCreateOfferIsValidBlockchain(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.TeamId = "2748779069441"
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.Hash()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	isValid, err := in.IsValidForBlockchain(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, false, isValid)
}
