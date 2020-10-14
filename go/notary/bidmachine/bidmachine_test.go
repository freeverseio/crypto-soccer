package bidmachine_test

import (
	"testing"
	"time"

	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/bidmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	"gotest.tools/assert"
)

func TestNotPayingAuction(t *testing.T) {
	auction := storage.Auction{State: storage.AuctionAssetFrozen}
	bid := storage.NewBid()
	shouldQueryMarketPay := true
	_, err := bidmachine.New(
		v1.NewMockMarketPay(),
		auction,
		bid,
		*bc.Contracts,
		bc.Owner,
		shouldQueryMarketPay,
	)
	assert.Error(t, err, "Auction is not in PAYING state")
}

func TestPayingAuction(t *testing.T) {
	auction := storage.Auction{State: storage.AuctionPaying}
	bid := storage.NewBid()
	shouldQueryMarketPay := true
	_, err := bidmachine.New(
		v1.NewMockMarketPay(),
		auction,
		bid,
		*bc.Contracts,
		bc.Owner,
		shouldQueryMarketPay,
	)
	assert.NilError(t, err)

	shouldQueryMarketPay = false
	_, err = bidmachine.New(
		v1.NewMockMarketPay(),
		auction,
		bid,
		*bc.Contracts,
		bc.Owner,
		shouldQueryMarketPay,
	)
	assert.NilError(t, err)
}

func TestFirstAlive(t *testing.T) {
	bid := bidmachine.FirstAlive(nil)
	if bid != nil {
		t.Fatalf("Wrong result: %v", bid)
	}
	bids := []storage.Bid{}
	bid = bidmachine.FirstAlive(bids)
	if bid != nil {
		t.Fatalf("Wrong result: %v", bid)
	}
	bids = []storage.Bid{storage.Bid{State: storage.BidFailed}}
	bid = bidmachine.FirstAlive(bids)
	if bid != nil {
		t.Fatalf("Wrong result: %v", bid)
	}
	bids = append(bids, storage.Bid{State: storage.BidAccepted, ExtraPrice: 10})
	bid = bidmachine.FirstAlive(bids)
	if *bid != bids[1] {
		t.Fatalf("Expected %v result: %v", bids[0], bid)
	}
	bids = append(bids, storage.Bid{State: storage.BidAccepted, ExtraPrice: 11})
	bid = bidmachine.FirstAlive(bids)
	if *bid != bids[2] {
		t.Fatalf("Wrong result: %v", bid)
	}
	bids = append(bids, storage.Bid{State: storage.BidPaying, ExtraPrice: 11})
	bid = bidmachine.FirstAlive(bids)
	if *bid != bids[3] {
		t.Fatalf("Wrong result: %v", bid)
	}
}

func TestExpiredBidNoTransit(t *testing.T) {
	auction := storage.Auction{State: storage.AuctionPaying}
	bid := &storage.Bid{State: storage.BidFailed}
	shouldQueryMarketPay := true
	machine, err := bidmachine.New(
		v1.NewMockMarketPay(),
		auction,
		bid,
		*bc.Contracts,
		bc.Owner,
		shouldQueryMarketPay,
	)
	assert.NilError(t, err)
	assert.NilError(t, machine.Process())
	assert.Equal(t, bid.State, storage.BidFailed)
}

func TestAcceptBidTransitToPaying(t *testing.T) {
	auction := storage.Auction{
		State:      storage.AuctionPaying,
		ValidUntil: 10,
	}
	bid := &storage.Bid{
		State:           storage.BidAccepted,
		PaymentDeadline: 3,
	}
	shouldQueryMarketPay := true
	machine, err := bidmachine.New(
		v1.NewMockMarketPay(),
		auction,
		bid,
		*bc.Contracts,
		bc.Owner,
		shouldQueryMarketPay,
	)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if bid.State != storage.BidPaying {
		t.Fatalf("Wrong state %v", bid.State)
	}
	// note: 172810 is almost exactly 48h, equal to the POST_AUCTION_TIME
	if bid.PaymentDeadline != 172810 {
		t.Fatalf("Wrong deadline %v", bid.PaymentDeadline)
	}
}

func TestBidPayingExpires(t *testing.T) {
	now := time.Now().Unix()
	auction := storage.Auction{
		Price:      3,
		State:      storage.AuctionPaying,
		ValidUntil: now - 3,
	}
	bid := &storage.Bid{
		State:           storage.BidPaying,
		PaymentDeadline: now - 1,
	}
	shouldQueryMarketPay := true
	machine, err := bidmachine.New(
		v1.NewMockMarketPay(),
		auction,
		bid,
		*bc.Contracts,
		bc.Owner,
		shouldQueryMarketPay,
	)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if bid.PaymentDeadline != now-1 {
		t.Fatalf("Wrong deadline %v", bid.PaymentDeadline)
	}
	if bid.State != storage.BidFailed {
		t.Fatalf("Wrong state %v", bid.State)
	}

	shouldQueryMarketPay = false
	machine, err = bidmachine.New(
		v1.NewMockMarketPay(),
		auction,
		bid,
		*bc.Contracts,
		bc.Owner,
		shouldQueryMarketPay,
	)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if bid.PaymentDeadline != now-1 {
		t.Fatalf("Wrong deadline %v", bid.PaymentDeadline)
	}
	if bid.State != storage.BidFailed {
		t.Fatalf("Wrong state %v", bid.State)
	}
}
