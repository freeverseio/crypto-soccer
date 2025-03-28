package auctionmachine_test

import (
	"testing"

	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestWithDrawableBySellerPendingValidate(t *testing.T) {
	market := v1.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	auction.State = storage.AuctionWithdrableBySeller
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}
	order.Status = "PENDING_VALIDATE"
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessWithdrawableBySeller(market))
	assert.Equal(t, m.State(), storage.AuctionValidation)
}

func TestWithDrawableBySellerPendingRelease(t *testing.T) {
	market := v1.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	auction.State = storage.AuctionWithdrableBySeller
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}
	order.Status = "PENDING_RELEASE"
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessWithdrawableBySeller(market))
	assert.Equal(t, m.State(), storage.AuctionWithdrableBySeller)
}

func TestWithDrawableBySellerReleased(t *testing.T) {
	market := v1.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	auction.State = storage.AuctionWithdrableBySeller
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}
	order.Status = "RELEASED"
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessWithdrawableBySeller(market))
	assert.Equal(t, m.State(), storage.AuctionWithdrableBySeller)
}
