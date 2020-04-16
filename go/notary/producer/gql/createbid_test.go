package gql_test

import (
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"

	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"gotest.tools/assert"
)

func TestCreateBidUnexistentAuction(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts)

	in := input.CreateBidInput{}
	_, err := r.CreateBid(struct{ Input input.CreateBidInput }{in})
	assert.Error(t, err, "")
}
