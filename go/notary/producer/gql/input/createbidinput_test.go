package input_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestCreateBidInputID(t *testing.T) {
	in := input.CreateBidInput{}
	assert.Equal(t, in.ID(), "f1534392279bddbf9d43dde8701cb5be14b82f76ec6607bf8d6ad557f60f304e")
}

func TestCreateBidInputHash(t *testing.T) {
	in := input.CreateBidInput{}
	auction := storage.NewAuction()
	auction.PlayerID = "2222"
	_, err := in.Hash(*bc.Contracts, *auction)
	assert.Error(t, err, "invalid teamId")
	in.TeamId = "10"
	hash, err := in.Hash(*bc.Contracts, *auction)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x2e873491a9bd59b403c457d71b4cdd119c3df561980c86ee0e2ea27d70e5d118")
}
