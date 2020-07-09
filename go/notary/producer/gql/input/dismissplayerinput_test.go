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
	assert.Equal(t, hash.Hex(), "0x26a63dd7a77ba6da621296c5433d235fa802b0eed629457ff3237b321f6db462")

	hash = helper.PrefixedHash(hash)
	assert.Equal(t, hash.Hex(), "0xa345906cc0144e72ba04ea426d34bd486000e51de093b4b1a106deafa21c3244")
}

func TestDismissPlayerValidSignature(t *testing.T) {
	msg := input.DismissPlayerInput{}
	msg.PlayerId = "123455"
	msg.ValidUntil = "5646456"
	msg.ReturnToAcademy = true
	msg.Signature = "2148732eeca5265898a5fe8dd3ba1c1af5b3d5b815fb23d9d6e383b376a2c91c694170ebd18b64b122905f82d3d6961a78a784b8966fcb350d51c6c5e7917d2d1b"

	isValid, err := msg.VerifySignature()
	assert.NilError(t, err)
	assert.Equal(t, isValid, true)
}

func TestDismissPlayerSignerAddress(t *testing.T) {
	msg := input.DismissPlayerInput{}
	msg.PlayerId = "123455"
	msg.ValidUntil = "5646456"
	msg.ReturnToAcademy = true
	msg.Signature = "2148732eeca5265898a5fe8dd3ba1c1af5b3d5b815fb23d9d6e383b376a2c91c694170ebd18b64b122905f82d3d6961a78a784b8966fcb350d51c6c5e7917d2d1b"

	address, err := msg.SignerAddress()
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), "0xb8CE9ab6943e0eCED004cDe8e3bBed6568B2Fa01")

	r, s, v, err := helper.RSV(msg.Signature)
	assert.Equal(t, hex.EncodeToString(r[:]), "2148732eeca5265898a5fe8dd3ba1c1af5b3d5b815fb23d9d6e383b376a2c91c")
	assert.Equal(t, hex.EncodeToString(s[:]), "694170ebd18b64b122905f82d3d6961a78a784b8966fcb350d51c6c5e7917d2d")
	assert.Equal(t, v, uint8(0x1b))
}
