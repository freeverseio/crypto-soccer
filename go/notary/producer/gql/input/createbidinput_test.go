package input_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/graph-gophers/graphql-go"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestCreateBidInputHash(t *testing.T) {
	auctioninput := input.CreateAuctionInput{}
	auctioninput.ValidUntil = "2000000000"
	auctioninput.OfferValidUntil = "1999999000"
	auctioninput.PlayerId = "274877906944"
	auctioninput.CurrencyId = 1
	auctioninput.Price = 41234
	auctioninput.Rnd = 42321
	auctionId, err := auctioninput.ID()
	assert.NilError(t, err)
	assert.Equal(t, string(auctionId), "24f45caee25883ab36cc32e2b152b94b8a05bdb086ac68784e3a4686f4d961e8")

	in := input.CreateBidInput{}
	in.AuctionId = auctionId
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = "274877906945"

	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x622481acfbf868f28f786b963f169b755390d06d88e5dfd08ebc9cdeab668dcf")
}

func TestCreateBidInputSign(t *testing.T) {
	auction := storage.NewAuction()
	auction.ValidUntil = int64(2000000000)
	auction.OfferValidUntil = int64(1999999000)
	auction.PlayerID = "274877906944"
	auction.CurrencyID = 1
	auction.Price = 41234
	auction.Rnd = 42321

	auctionId, err := auction.ComputeID()
	assert.Equal(t, auctionId, "24f45caee25883ab36cc32e2b152b94b8a05bdb086ac68784e3a4686f4d961e8")

	in := input.CreateBidInput{}
	in.AuctionId = graphql.ID(auctionId)
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = "274877906945"
	in.Signature = "4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b"

	isValid, err := in.IsSignerOwner(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, isValid, false)
}
