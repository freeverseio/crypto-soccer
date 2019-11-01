package bidmachine_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/bidmachine"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

func TestNotPayingAuction(t *testing.T) {
	auction := storage.Auction{State: storage.AUCTION_ASSET_FROZEN}
	bids := []storage.Bid{}
	_, err := bidmachine.New(auction, bids)
	if err == nil {
		t.Fatalf("Accepting %v auction", auction.State)
	}
}

func TestPayingAuction(t *testing.T) {
	auction := storage.Auction{State: storage.AUCTION_PAYING}
	bids := []storage.Bid{}
	_, err := bidmachine.New(auction, bids)
	if err != nil {
		t.Fatalf("Not accepting %v auction", auction.State)
	}
}

func TestExpiredBidNoTransit(t *testing.T) {
	auction := storage.Auction{State: storage.AUCTION_PAYING}
	bids := []storage.Bid{storage.Bid{State: storage.BID_EXPIRED}}
	machine, err := bidmachine.New(auction, bids)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if machine.Bids[0].State != storage.BID_EXPIRED {
		t.Fatalf("Wrong state %v", machine.Bids[0].State)
	}
}

func TestAcceptBidTransitToPaying(t *testing.T) {
	auction := storage.Auction{State: storage.AUCTION_PAYING}
	bids := []storage.Bid{storage.Bid{State: storage.BID_ACCEPTED}}
	machine, err := bidmachine.New(auction, bids)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if machine.Bids[0].State != storage.BID_PAYING {
		t.Fatalf("Wrong state %v", machine.Bids[0].State)
	}
}
