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
