package input_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
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
	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x6c916b19d781786ba2a3e90fbb3c772df29c8917d1b0486f92ff40234e5cdcdc")
}

func TestCreateOfferInputID(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	in.TeamId = "20"
	id, err := in.ID(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, string(id), "6c916b19d781786ba2a3e90fbb3c772df29c8917d1b0486f92ff40234e5cdcdc")
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
	valid, err := in.VerifySignature(*bc.Contracts)
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
	address, err := in.SignerAddress(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0xA39Deb9bda8CaED16Ae57667a3F8BEcD6970F559")
}

func TestCreateOfferIsSignerOwner(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "27487790694"
	in.CurrencyId = 1
	in.TeamId = "20"
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	isOwner, err := in.IsSignerOwner(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, isOwner, false)
}

func TestCreateOfferGetOwner(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.TeamId = "20"
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	owner, err := in.GetOwner(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, crypto.PubkeyToAddress(bc.Owner.PublicKey).Hex(), owner.Hex())
}
