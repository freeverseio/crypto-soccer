package input_test

import (
	"encoding/hex"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestSubmitPlayStorePlayerPurchaseInputHash(t *testing.T) {
	in := input.SubmitPlayStorePlayerPurchaseInput{}

	hash, err := in.Hash()
	assert.Error(t, err, "Invalid TeamId")

	in.TeamId = "3"
	hash, err = in.Hash()
	assert.Error(t, err, "Invalid PlayerId")

	in.PlayerId = "5"
	in.Receipt = "korpimulxmslxissnschtkdb"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0xaea70cc97f358a6f00d5bed69952c06846423d450ad73fd7d6f385905dd18619")
}

func TestSubmitPlayStorePlayerPurchaseInputSignature(t *testing.T) {
	in := input.SubmitPlayStorePlayerPurchaseInput{}
	in.TeamId = "274877906944"
	in.PlayerId = "274877906944"
	in.Receipt = "com.freeverse.phoenix"

	hash, err := in.Hash()
	assert.NilError(t, err)
	hash = helper.PrefixedHash(hash)
	sign, err := helper.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "872caa6f47258638007034bbd4d84002c2e17b9e9a7607db25bb8164d22d7e751f83f084442f98a0c4b647b15286209ea9eb52807c27bb5777af13aaa6110fc81b")
}
