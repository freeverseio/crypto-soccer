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
	service := postgres.NewStorageService(db)
	assert.NilError(t, service.Begin())
	defer service.Rollback()

	in := input.CreateBidInput{}
	assert.Error(t, consumer.CreateBid(service, in), "No auction for bid {  0 0 }")
}

func TestCreateBid(t *testing.T) {
	service := postgres.NewStorageService(db)
	assert.NilError(t, service.Begin())
	defer service.Rollback()

	auction := storage.NewAuction()
	auction.ID = "3"
	assert.NilError(t, service.AuctionInsert(*auction))

	in := input.CreateBidInput{}
	in.AuctionId = graphql.ID(auction.ID)
	assert.NilError(t, consumer.CreateBid(service, in))
}
