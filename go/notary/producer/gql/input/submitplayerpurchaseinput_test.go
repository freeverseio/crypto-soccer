package input_test

import (
	"encoding/hex"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestSubmitPlayerPurchaseInputHash(t *testing.T) {
	in := input.SubmitPlayerPurchaseInput{}

	hash, err := in.Hash()
	assert.Error(t, err, "Invalid TeamId")

	in.TeamId = "3"
	hash, err = in.Hash()
	assert.Error(t, err, "Invalid PlayerId")

	in.PlayerId = "5"
	in.PurchaseId = "korpimulxmslxissnschtkdb"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x03d0fa8eac063046346e9694c5d38bdb1b180017c79aa1e70384becbd4158f4d")
}

func TestSubmitPlayerPurchaseInputSignature(t *testing.T) {
	in := input.SubmitPlayerPurchaseInput{}
	in.TeamId = "274877906944"
	in.PlayerId = "274877906944"
	in.PurchaseId = "korpimulxmslxissnschtkdb"

	hash, err := in.Hash()
	assert.NilError(t, err)
	hash = input.PrefixedHash(hash)
	sign, err := input.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "91366deb26195ac3b15b9e6fff99d425b2cc0d15d44dc8ee0377779400f92c4358a57754053facbe724e8a536e240b278cd651c756c46978eaebafc47767fd781b")
}
