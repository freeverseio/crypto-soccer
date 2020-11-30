package input_test

import (
	"encoding/hex"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestSubmitAuctionPassPlayStorePurchaseInputHash(t *testing.T) {
	in := input.SubmitAuctionPassPlayStorePurchaseInput{}

	hash, err := in.Hash()
	assert.Error(t, err, "Invalid TeamId")

	in.TeamId = "3"

	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x6183a04476692ce45a4414be1bb7d2df6bd59a2adfb90b95fdcabca89bfa0295")
}

func TestSubmitAuctionPassPlayStorePurchaseInputSignature(t *testing.T) {
	in := input.SubmitAuctionPassPlayStorePurchaseInput{}
	in.TeamId = "274877906944"
	in.Receipt = "com.freeverse.phoenix"

	hash, err := in.Hash()
	assert.NilError(t, err)
	hash = helper.PrefixedHash(hash)
	sign, err := helper.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "717fe1b4b63e4450244e6d89827be4cd22e7429c87cd52bd5dfed062b1cb967d69d3fc5b83800da4ff4c4cd4b5b4ca800458903f3005dd43c80a9c08abad0aba1b")
}
