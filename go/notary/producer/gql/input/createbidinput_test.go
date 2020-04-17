package input_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestCreateBidInputHash(t *testing.T) {
	auction := input.CreateAuctionInput{}
	auction.ValidUntil = "2000000000"
	auction.PlayerId = "274877906944"
	auction.CurrencyId = 1
	auction.Price = 41234
	auction.Rnd = 42321
	auctionId, err := auction.ID()
	assert.NilError(t, err)
	assert.Equal(t, string(auctionId), "55d0b23ce4ce7530aa71b177b169ca4bf52dec4866ffbf37fa84fd0146a5f36a")

	in := input.CreateBidInput{}
	in.Auction = auctionId
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = "274877906945"

	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0xc0ad1683b9afe071d698763b7143e7cff7bcc661c7074497d870964dd58d9976")
}

func TestCreateBidInputSign(t *testing.T) {
	auction := storage.NewAuction()
	auction.ValidUntil = int64(2000000000)
	auction.PlayerID = "274877906944"
	auction.CurrencyID = 1
	auction.Price = 41234
	auction.Rnd = 42321

	in := input.CreateBidInput{}
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = "274877906945"
	in.Signature = "4fe5772189b4e448e528257f6b32b3ebc90ed8f52fc7c9b04594d86adb74875147f62c6d83b8555c63d622b2248bb6846c75912a684490a68de46ede201ecf0f1b"

	isValid, err := in.VerifySignature(*bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, isValid)
}
