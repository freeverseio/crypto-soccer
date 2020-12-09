package gql_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/relay/producer/gql"
	"gotest.tools/assert"
)

func TestConsumePromoHash(t *testing.T) {
	in := gql.ConsumePromoInput{}
	in.TeamId = "64645"
	in.PlayerId = "333"
	hash, err := in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x55314d0664b112c63cce9758604e0656eca3ce7958536f1b175374f258e49ff9")
	in.PlayerId = "2"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x77ea2a719178d7a661ebfa08532c933b85ed56c55081d4bf6b7b08cd25978d4d")
}

func TestConsumePromoGetSigner(t *testing.T) {
	in := gql.ConsumePromoInput{}
	in.TeamId = "64645"
	in.PlayerId = "333"

	hash, err := in.Hash()
	assert.NilError(t, err)

	pvc, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	sign, err := signer.Sign(hash.Bytes(), pvc)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "06ef681489c0c9a2d8bf6b9f862acea1c62354e447ecb55ce9ee837bf3c1bd09023f6ea5395f7f9eb70b7000756f03a424cd12aa76ee97f6030759bb7acc66231c")

	address, err := in.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
}
