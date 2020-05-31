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
	hash = input.PrefixedHash(hash)
	sign, err := input.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "82b6568d3e792df067a07ca67316b916de3064ef0cdabcbf25a59e5e9745caa328ae510bd2a62a92e2f9710aa38798a0a7e7f47b0632bf08fa4c7abd52e5c0a11b")

	isValid, err := input.VerifySignature(hash, sign)
	assert.NilError(t, err)
	assert.Assert(t, isValid)

	sender, err := in.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, sender.Hex(), "0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7")
}

func TestGeneratePlayerIdsSignature2(t *testing.T) {
	in := input.GetWorldPlayersInput{}
	in.TeamId = "274877906944"

	hash, err := in.Hash()
	assert.NilError(t, err)
	hash = input.PrefixedHash(hash)
	sign, err := input.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "a67621b4763db406f404c4a600ce0e79ee50147c209e85d2f146f0d760c0a1ac2a213a06f702995cee279af1f588b55c9fa462b2e6a9502d25cede77ec690ced1c")
}

func TestGeneratePlayerIdsSignerIsOwnerOfTeam(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	assert.NilError(t, err)
	in := input.GetWorldPlayersInput{}
	in.TeamId = "4"
	hash, err := in.Hash()
	assert.NilError(t, err)
	hash = input.PrefixedHash(hash)
	sign, err := input.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(sign)
	assert.Equal(t, in.Signature, "82b6568d3e792df067a07ca67316b916de3064ef0cdabcbf25a59e5e9745caa328ae510bd2a62a92e2f9710aa38798a0a7e7f47b0632bf08fa4c7abd52e5c0a11b")

	isOwner, err := in.IsSignerOwner(*bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, !isOwner)
}
