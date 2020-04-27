package auctionmachine_test

import (
	"testing"

	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestWithDrawableBySellerPendingValidate(t *testing.T) {
	market := marketpay.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	auction.State = storage.AuctionWithdrableBySeller
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}
	order.Status = "PENDING_VALIDATE"
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessWithdrawableBySeller(market))
	assert.Equal(t, m.State(), storage.AuctionWithdrableBySeller)
}

func TestWithDrawableBySellerPendingRelease(t *testing.T) {
	market := marketpay.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	auction.State = storage.AuctionWithdrableBySeller
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}
	order.Status = "PENDING_RELEASE"
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessWithdrawableBySeller(market))
	assert.Equal(t, m.State(), storage.AuctionEnded)
}

func TestWithDrawableBySellerReleased(t *testing.T) {
	market := marketpay.NewMockMarketPay()
	order, err := market.CreateOrder("Bellaciao", "1000.0")
	assert.NilError(t, err)
	auction := storage.NewAuction()
	auction.State = storage.AuctionWithdrableBySeller
	bid := storage.NewBid()
	bid.State = storage.BidPaid
	bid.PaymentID = order.TrusteeShortlink.Hash
	bids := []storage.Bid{*bid}
	order.Status = "RELEASED"
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessWithdrawableBySeller(market))
	assert.Equal(t, m.State(), storage.AuctionEnded)
}
