package bidmachine_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/bidmachine"
	"github.com/freeverseio/crypto-soccer/go/testutils"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

func TestNotPayingAuction(t *testing.T) {
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{State: storage.AUCTION_ASSET_FROZEN}
	bid := storage.Bid{}
	_, err = bidmachine.New(
		auction,
		bid,
		bc.Market,
		bc.Owner,
		bc.Client,
	)
	if err == nil {
		t.Fatalf("Accepting %v auction", auction.State)
	}
}

func TestPayingAuction(t *testing.T) {
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{State: storage.AUCTION_PAYING}
	bid := storage.Bid{}
	_, err = bidmachine.New(
		auction,
		bid,
		bc.Market,
		bc.Owner,
		bc.Client,
	)
	if err != nil {
		t.Fatalf("Not accepting %v auction", auction.State)
	}
}

func TestFirstAlive(t *testing.T) {
	idx := bidmachine.IndexFirstAlive(nil)
	if idx != -1 {
		t.Fatalf("Wrong result: %v", idx)
	}
	bids := []storage.Bid{}
	idx = bidmachine.IndexFirstAlive(bids)
	if idx != -1 {
		t.Fatalf("Wrong result: %v", idx)
	}
	bids = []storage.Bid{storage.Bid{State: storage.BID_FAILED_TO_PAY}}
	idx = bidmachine.IndexFirstAlive(bids)
	if idx != -1 {
		t.Fatalf("Wrong result: %v", idx)
	}
	bids = append(bids, storage.Bid{State: storage.BID_ACCEPTED, ExtraPrice: 10})
	idx = bidmachine.IndexFirstAlive(bids)
	if idx != 1 {
		t.Fatalf("Wrong result: %v", idx)
	}
	bids = append(bids, storage.Bid{State: storage.BID_ACCEPTED, ExtraPrice: 11})
	idx = bidmachine.IndexFirstAlive(bids)
	if idx != 2 {
		t.Fatalf("Wrong result: %v", idx)
	}
	bids = append(bids, storage.Bid{State: storage.BID_PAYING, ExtraPrice: 11})
	idx = bidmachine.IndexFirstAlive(bids)
	if idx != 3 {
		t.Fatalf("Wrong result: %v", idx)
	}
}

func TestExpiredBidNoTransit(t *testing.T) {
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{State: storage.AUCTION_PAYING}
	bid := storage.Bid{State: storage.BID_FAILED_TO_PAY}
	machine, err := bidmachine.New(
		auction,
		bid,
		bc.Market,
		bc.Owner,
		bc.Client,
	)
	if err != nil {
		t.Fatal(err)
	}
	bid, err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if bid.State != storage.BID_FAILED_TO_PAY {
		t.Fatalf("Wrong state %v", bid.State)
	}
}

func TestAcceptBidTransitToPaying(t *testing.T) {
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{State: storage.AUCTION_PAYING}
	bid := storage.Bid{State: storage.BID_ACCEPTED}
	machine, err := bidmachine.New(
		auction,
		bid,
		bc.Market,
		bc.Owner,
		bc.Client,
	)
	if err != nil {
		t.Fatal(err)
	}
	bid, err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if bid.State != storage.BID_PAYING {
		t.Fatalf("Wrong state %v", bid.State)
	}
}
