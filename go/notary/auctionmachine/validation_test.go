package auctionmachine_test

import (
	"testing"

	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestValidationAuctionInvalidState(t *testing.T) {
	market := marketpay.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}

	auction.State = storage.AuctionWithdrableBySeller
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner)
	assert.NilError(t, err)
	assert.Error(t, m.ProcessValidation(market), "Wrong state withadrable_by_seller")

	auction.State = storage.AuctionAssetFrozen
	m, err = auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner)
	assert.NilError(t, err)
	assert.Error(t, m.ProcessValidation(market), "Wrong state asset_frozen")
}

func TestValidationAuctionValidOrderInvalidState(t *testing.T) {
	market := marketpay.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}

	order.Status = "DRAFT"
	auction.State = storage.AuctionValidation
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessValidation(market))
	assert.Equal(t, m.State(), storage.AuctionValidation)
}

func TestValidationAuctionValidOrderPendingRelease(t *testing.T) {
	market := marketpay.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}

	order.Status = "PENDING_RELEASE"
	auction.State = storage.AuctionValidation
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessValidation(market))
	assert.Equal(t, m.State(), storage.AuctionValidation)
}

func TestValidationAuctionValidOrderReleased(t *testing.T) {
	market := marketpay.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}

	order.Status = "RELEASED"
	auction.State = storage.AuctionValidation
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessValidation(market))
	assert.Equal(t, m.State(), storage.AuctionEnded)
}
