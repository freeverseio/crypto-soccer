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
	assert.Equal(t, hash.Hex(), "0x09ee27431fe606a3c506131fbd025a440d4318c40f1500ee0dea757cd6c67ccb")
}
func TestCancelAllOffersBySellerGetSigner(t *testing.T) {
	in := input.CancelAllOffersBySellerInput{}
	in.PlayerId = "274877906944"

	hash, err := in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x09ee27431fe606a3c506131fbd025a440d4318c40f1500ee0dea757cd6c67ccb")

	pvc, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	sign, err := signer.Sign(hash.Bytes(), pvc)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "a67621b4763db406f404c4a600ce0e79ee50147c209e85d2f146f0d760c0a1ac2a213a06f702995cee279af1f588b55c9fa462b2e6a9502d25cede77ec690ced1c")

	address, err := in.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
}
