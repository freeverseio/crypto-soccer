package input_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

func TestCancelAllOffersBySellerHash(t *testing.T) {
	in := input.CancelAllOffersBySellerInput{}
	in.PlayerId = ""
	hash, err := in.Hash()
	assert.Error(t, err, "invalid playerId")
	in.PlayerId = "274877906944"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x978c23f58e78cae5d83b11b244d04c67bd802e70f7fda94bf940379a9713812e")
}
func TestCancelAllOffersBySellerGetSigner(t *testing.T) {
	in := input.CancelAllOffersBySellerInput{}
	in.PlayerId = "274877906944"

	hash, err := in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x978c23f58e78cae5d83b11b244d04c67bd802e70f7fda94bf940379a9713812e")

	pvc, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	sign, err := signer.Sign(hash.Bytes(), pvc)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "8b641b11779fd458998b04fc90c35bb717237e72a8fa6ce045b1a38ad610a1f16b51d05a0521a65553167b9b4b8e4a5286648a727bc9b36a496d64bcaa5658a31b")

	address, err := in.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
}
