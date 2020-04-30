package input_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestGeneratePlayerIdsHash(t *testing.T) {
	in := input.GetWorldPlayersInput{}
	in.TeamId = "4"
	hash, err := in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b")

	in.Signature = "dsdsd"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b")

	in.TeamId = "5"
	hash, err = in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x036b6384b5eca791c62761152d0c79bb0604c104a5fb6f4eb0703f3154bb3db0")
}

func TestGeneratePlayerIdsSignature(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)

	in := input.GetWorldPlayersInput{}
	in.TeamId = "4"

	hash, err := in.Hash()
	assert.NilError(t, err)
	sign, err := input.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "724d892b43b6640b3baa00236f6920caa217596b313f999d32e8737019e65e04684e74d532b71f34c73c62c5783387565c5ac3e03cfb4b9b950c0859899339481c")

	isValid, err := input.VerifySignature(hash, sign)
	assert.NilError(t, err)
	assert.Assert(t, isValid)

	sender, err := input.AddressFromSignature(hash, sign)
	assert.NilError(t, err)
	assert.Equal(t, sender.Hex(), "0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7")
}
