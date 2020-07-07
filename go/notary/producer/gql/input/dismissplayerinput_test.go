package input_test

import (
	"encoding/hex"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestDismissPlayerInputHash(t *testing.T) {
	msg := input.DismissPlayerInput{}
	msg.PlayerId = "123455"
	msg.ValidUntil = "5646456"
	msg.ReturnToAcademy = true

	hash, err := msg.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0xdaa54d5f301de697c9dab3416a20ac9d4d7890f92f01c1340c18152cb93c5ca0")

	hash = helper.PrefixedHash(hash)
	assert.Equal(t, hash.Hex(), "0xfde13974fe686e362246c490ef0280e35928eca805695c7bfd89c7f2bc74b39b")
}

func TestDismissPlayerValidSignature(t *testing.T) {
	msg := input.DismissPlayerInput{}
	msg.PlayerId = "123455"
	msg.ValidUntil = "5646456"
	msg.ReturnToAcademy = true
	msg.Signature = "0f13e4028d911bbf7e305267d593c6b67888030032e73f94a5cf8af204567ab629848e9290568aa5d19c1b7a4761a20ed4059072aacd79bde56e1b52c17a21311b"

	isValid, err := msg.VerifySignature()
	assert.NilError(t, err)
	assert.Equal(t, isValid, true)
}

func TestDismissPlayerSignerAddress(t *testing.T) {
	msg := input.DismissPlayerInput{}
	msg.PlayerId = "123455"
	msg.ValidUntil = "5646456"
	msg.ReturnToAcademy = true
	msg.Signature = "0f13e4028d911bbf7e305267d593c6b67888030032e73f94a5cf8af204567ab629848e9290568aa5d19c1b7a4761a20ed4059072aacd79bde56e1b52c17a21311b"

	address, err := msg.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0xb8CE9ab6943e0eCED004cDe8e3bBed6568B2Fa01")

	r, s, v, err := helper.RSV(msg.Signature)
	assert.Equal(t, hex.EncodeToString(r[:]), "0f13e4028d911bbf7e305267d593c6b67888030032e73f94a5cf8af204567ab6")
	assert.Equal(t, hex.EncodeToString(s[:]), "29848e9290568aa5d19c1b7a4761a20ed4059072aacd79bde56e1b52c17a2131")
	assert.Equal(t, v, uint8(0x1b))
}
