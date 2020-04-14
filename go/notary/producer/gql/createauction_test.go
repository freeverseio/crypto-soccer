package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"gotest.tools/assert"
)

func TestCreateAuctionReturnTheSignature(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch)

	in := input.CreateAuctionInput{}
	in.Signature = "534523re32"

	id, err := r.CreateAuction(struct{ Input input.CreateAuctionInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, string(id), in.Signature)
}
