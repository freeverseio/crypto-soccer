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
	in.PurchaseToken = "korpimulxmslxissnschtkdb"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x0d64515f0387f90deddb55c0bb8581719a08090575cc0dad3ba129cf4c0dc53e")

	in.ProductId = "player1000"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x68fbd817653deba84c34cde2d67821141e38b78abe3f04b654b675f6b383e373")

	in.PackageName = "pippo"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x9465358509150f759185df7c9b9100a951d47643697605bc4822ca0bdd5ce018")
}

func TestSubmitPlayStorePlayerPurchaseInputSignature(t *testing.T) {
	in := input.SubmitPlayStorePlayerPurchaseInput{}
	in.TeamId = "274877906944"
	in.PlayerId = "274877906944"
	in.PackageName = "com.freeverse.phoenix"
	in.ProductId = "coinpack_45"
	in.PurchaseToken = "korpimulxmslxissnschtkdb"

	hash, err := in.Hash()
	assert.NilError(t, err)
	hash = helper.PrefixedHash(hash)
	sign, err := helper.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "11baa07b721ac1991c266973c70a712bfb261da9133d4f5c38e9019e8575efe56ca32f13664949251b6cb0620389f96662849e00efbb5d5f347e211379f9e4151c")
}
