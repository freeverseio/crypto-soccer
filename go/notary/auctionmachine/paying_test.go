package auctionmachine_test

import (
	"testing"
	"time"

	"github.com/freeverseio/crypto-soccer/go/marketpay"
	"github.com/freeverseio/crypto-soccer/go/marketpay/mock"
	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestPaying(t *testing.T) {
	auction := storage.NewAuction()
	auction.ID = "f1d4501c5158a9018b1618ec4d471c66b663d8f6bffb6e70a0c6584f5c1ea94a"
	auction.ValidUntil = time.Now().Unix() + 100
	auction.PlayerID = "274877906944"
	auction.CurrencyID = 1
	auction.Price = 41234
	auction.Rnd = 4232
	auction.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e"
	auction.Signature = "381bf58829e11790830eab9924b123d1dbe96dd37b10112729d9d32d476c8d5762598042bb5d5fd63f668455aa3a2ce4e2632241865c26ababa231ad212b5f151b"

	auction.State = storage.AuctionStarted
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, nil, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.Error(t, m.ProcessPaying(v1.NewMockMarketPay()), "auction[f1d4501c5158a9018b1618ec4d471c66b663d8f6bffb6e70a0c6584f5c1ea94a|started] is not in state paying")
}

func TestPayingWithOffer(t *testing.T) {
	auction := storage.NewAuction()
	auction.ID = "f1d4501c5158a9018b1618ec4d471c66b663d8f6bffb6e70a0c6584f5c1ea94a"
	auction.ValidUntil = time.Now().Unix() + 100
	auction.PlayerID = "274877906944"
	auction.CurrencyID = 1
	auction.Price = 41234
	auction.Rnd = 4232
	auction.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e"
	auction.Signature = "381bf58829e11790830eab9924b123d1dbe96dd37b10112729d9d32d476c8d5762598042bb5d5fd63f668455aa3a2ce4e2632241865c26ababa231ad212b5f151b"

	auction.State = storage.AuctionStarted
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, nil, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.Error(t, m.ProcessPaying(v1.NewMockMarketPay()), "auction[f1d4501c5158a9018b1618ec4d471c66b663d8f6bffb6e70a0c6584f5c1ea94a|started] is not in state paying")
}

func TestPayingNilBids(t *testing.T) {
	auction := storage.NewAuction()
	auction.State = storage.AuctionPaying
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, nil, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessPaying(v1.NewMockMarketPay()))
	assert.Equal(t, m.State(), storage.AuctionFailed)
	assert.Equal(t, m.StateExtra(), "No available healty bid")
}

func TestPayingNoBidsAvailable(t *testing.T) {
	auction := storage.NewAuction()
	auction.State = storage.AuctionPaying
	bid := storage.NewBid()
	bid.State = storage.BidFailed
	bids := []storage.Bid{*bid}
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessPaying(v1.NewMockMarketPay()))
	assert.Equal(t, m.State(), storage.AuctionFailed)
	assert.Equal(t, m.StateExtra(), "No available healty bid")
}

func TestPayingWithBid(t *testing.T) {
	auction := storage.NewAuction()
	auction.State = storage.AuctionPaying
	bid := storage.NewBid()
	bid.State = storage.BidAccepted
	bids := []storage.Bid{*bid}
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessPaying(v1.NewMockMarketPay()))
	assert.Equal(t, m.State(), storage.AuctionPaying)
}

func TestPayingFailing(t *testing.T) {
	service := mock.MarketPayMock{
		CreateOrderFunc: func(name string, value string) (*marketpay.Order, error) {
			return &marketpay.Order{
				Status: "PUBLISHED",
			}, nil
		},
		GetOrderFunc: func(hash string) (*marketpay.Order, error) {
			return &marketpay.Order{}, nil
		},
		IsPaidFunc: func(order marketpay.Order) bool {
			return true
		},
	}

	auction := storage.NewAuction()
	auction.ID = "ciao"
	auction.State = storage.AuctionPaying
	auction.PlayerID = "43253543"
	auction.CurrencyID = 1
	auction.Price = 20
	auction.Rnd = 44
	auction.Signature = "381bf58829e11790830eab9924b123d1dbe96dd37b10112729d9d32d476c8d5762598042bb5d5fd63f668455aa3a2ce4e2632241865c26ababa231ad212b5f151b"
	bid := storage.NewBid()
	bid.AuctionID = auction.ID
	bid.State = storage.BidPaying
	bid.PaymentDeadline = time.Now().Unix() + 10000
	bid.TeamID = "4343"
	bid.Signature = "381bf58829e11790830eab9924b123d1dbe96dd37b10112729d9d32d476c8d5762598042bb5d5fd63f668455aa3a2ce4e2632241865c26ababa231ad212b5f151b"
	bids := []storage.Bid{*bid}
	shouldQueryMarketPay := true
	m, err := auctionmachine.New(*auction, bids, *bc.Contracts, bc.Owner, shouldQueryMarketPay)
	assert.NilError(t, err)
	assert.NilError(t, m.ProcessPaying(&service))
	assert.Equal(t, m.State(), storage.AuctionWithdrableByBuyer)
}
