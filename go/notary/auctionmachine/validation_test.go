package auctionmachine_test

import (
	"testing"

	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestValidationAuctionInvalidState(t *testing.T) {
	market := v1.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}

	auction.State = storage.AuctionWithdrableBySeller
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.Error(t, m.ProcessValidation(market), "auction[|withadrable_by_seller] is not in state validation")

	auction.State = storage.AuctionAssetFrozen
	m, err = auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.Error(t, m.ProcessValidation(market), "auction[|asset_frozen] is not in state validation")
}

func TestValidationAuctionValidOrderInvalidState(t *testing.T) {
	market := v1.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}

	order.Status = "DRAFT"
	auction.State = storage.AuctionValidation
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessValidation(market))
	assert.Equal(t, m.State(), storage.AuctionValidation)
}

func TestValidationAuctionValidOrderPendingRelease(t *testing.T) {
	market := v1.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}

	order.Status = "PENDING_RELEASE"
	auction.State = storage.AuctionValidation
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessValidation(market))
	assert.Equal(t, m.State(), storage.AuctionValidation)
}

func TestValidationAuctionValidOrderReleased(t *testing.T) {
	market := v1.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}

	order.Status = "RELEASED"
	auction.State = storage.AuctionValidation
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessValidation(market))
	assert.Equal(t, m.State(), storage.AuctionEnded)
}
