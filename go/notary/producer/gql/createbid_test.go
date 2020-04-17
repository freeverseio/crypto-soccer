package gql_test

import (
	"strconv"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"

	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestCreateBidUnexistentAuction(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, db)

	in := input.CreateBidInput{}
	_, err := r.CreateBid(struct{ Input input.CreateBidInput }{in})
	assert.Error(t, err, "unexistent auction")
}

func TestCreateBid(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, db)

	auction := input.CreateAuctionInput{}
	auction.ValidUntil = "2000000000"
	auction.PlayerId = "274877906944"
	auction.CurrencyId = 1
	auction.Price = 41234
	auction.Rnd = 42321
	auctionId := auction.ID()

	auctionS := storage.NewAuction()
	auctionS.ID = auctionId
	validUntil, _ := strconv.ParseInt(auction.ValidUntil, 10, 64)
	auctionS.ValidUntil = validUntil
	auctionS.PlayerID = auction.PlayerId
	auctionS.CurrencyID = int(auction.CurrencyId)
	auctionS.Price = int(auction.Price)
	auctionS.Rnd = int(auction.Rnd)
	assert.NilError(t, auctionS.Insert(tx))

	in := input.CreateBidInput{}
	in.Auction = graphql.ID(auctionId)
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = "274877906945"

	id, err := r.CreateBid(struct{ Input input.CreateBidInput }{in})
	assert.NilError(t, err)
	assert.Equal(t, string(id), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
}
