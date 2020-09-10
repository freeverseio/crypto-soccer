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
	assert.Equal(t, string(auctionId), "aa5d5b3de11b03fe9def7911646e0661ce335d423fc5c740b3db49b11b9f7604")

	in := input.CreateBidInput{}
	in.AuctionId = auctionId
	in.ExtraPrice = 332
	in.Rnd = 1243523
	in.TeamId = "274877906945"

	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x0dd3cc175ce499ed8845a45af76e875a54e036b638ba96e726143ff567c2b280")
}

func TestCreateBidInputSign(t *testing.T) {
	auction := storage.NewAuction()
	auction.ValidUntil = uint32(2000000000)
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
