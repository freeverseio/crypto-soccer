package consumer_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/graph-gophers/graphql-go"
	"gotest.tools/assert"
)

func TestCreateBidWithNoAuction(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	in := input.CreateBidInput{}
	assert.Error(t, consumer.CreateBid(tx, in), "No auction for bid {  0 0 }")
}

func TestCreateBid(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	service := postgres.NewAuctionService(tx)

	auction := storage.NewAuction()
	auction.ID = "3"
	assert.NilError(t, service.Insert(*auction))

	in := input.CreateBidInput{}
	in.AuctionId = graphql.ID(auction.ID)
	assert.NilError(t, consumer.CreateBid(tx, in))
}
