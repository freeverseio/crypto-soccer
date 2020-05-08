package input_test

import (
	"encoding/hex"
	"testing"

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
	in.PurchaseToken = "korpimulxmslxissnschtkdb"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x273b5c8787548a907a75b13725cca7f842b9f881906686846c0d4430422ebca6")

	in.ProductId = "player1000"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0xab052e01c0cfcabb8b121c2c46d6cb36fa0e35005519000eaed33e60f47b0574")
}

func TestSubmitPlayStorePlayerPurchaseInputSignature(t *testing.T) {
	in := input.SubmitPlayStorePlayerPurchaseInput{}
	in.TeamId = "274877906944"
	in.PlayerId = "274877906944"
	in.ProductId = "coinpack_45"
	in.PurchaseToken = "korpimulxmslxissnschtkdb"

	hash, err := in.Hash()
	assert.NilError(t, err)
	hash = input.PrefixedHash(hash)
	sign, err := input.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "a4f900f8923554e48acb7977f78b8569e3bc6e2d80ffbbd040e79e6bf602a1491c43995ee4ab2e77dbfae24f4a870c95a4db11d3276ed423cd49b0bc18c45a091c")
}
