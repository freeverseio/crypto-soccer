package input_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestCreateAuctionInputHash(t *testing.T) {
	in := input.CreateAuctionInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	hash, err := in.Hash()
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0xc50d978b8a838b6c437a162a94c715f95e92e11fe680cf0f1caf054ad78cd796")
}
