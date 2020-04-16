package input_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestCreateBidInputID(t *testing.T) {
	in := input.CreateBidInput{}
	assert.Equal(t, in.ID(), "f1534392279bddbf9d43dde8701cb5be14b82f76ec6607bf8d6ad557f60f304e")
}

func TestCreateBidInputHash(t *testing.T) {
	in := input.CreateBidInput{}
	_, err := in.Hash()
	assert.Error(t, err, "invalid teamId")
	in.TeamId = "10"
	hash, err := in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x0")
}
