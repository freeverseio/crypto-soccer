package gql_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"gotest.tools/assert"
)

func TestCreateAuctionReturnTheSignature(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch)

	input := gql.CreateAuctionInput{}
	input.Signature = "534523re32"

	id, err := r.CreateAuction(struct{ Input gql.CreateAuctionInput }{input})
	assert.NilError(t, err)
	assert.Equal(t, string(id), input.Signature)
}
