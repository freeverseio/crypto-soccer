package gql_test

import (
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"

	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"gotest.tools/assert"
)

func TestCreateBidUnexistentAuction(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, db)

	in := input.CreateBidInput{}
	in.TeamId = "10"
	_, err := r.CreateBid(struct{ Input input.CreateBidInput }{in})
	assert.Error(t, err, "signature must be 65 bytes long")
}

func TestCreateBid(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, db)

	auction := input.CreateAuctionInput{}
	auction.ValidUntil = "2000000000"
	auction.PlayerId = "274877906944"
	auction.CurrencyId = 1
	auction.Price = 41234
	auction.Rnd = 42321
	auctionId, err := auction.ID()
	assert.NilError(t, err)

	in := input.CreateBidInput{}
	in.AuctionId = auctionId
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = "274877906945"
	in.Signature = "4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1c"

	id, err := r.CreateBid(struct{ Input input.CreateBidInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, string(id), "c0ad1683b9afe071d698763b7143e7cff7bcc661c7074497d870964dd58d9976")
}
