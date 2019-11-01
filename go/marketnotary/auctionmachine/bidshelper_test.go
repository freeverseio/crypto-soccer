package auctionmachine_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/auctionmachine"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
)

func TestOrderByDEscExtraPrice(t *testing.T) {
	bids := []storage.Bid{}
	bids = auctionmachine.OrderByDescExtraPrice(bids)
	if len(bids) != 0 {
		t.Fatal("Not 0 lenght")
	}

	bids = []storage.Bid{
		storage.Bid{
			ExtraPrice: 3,
		},
		storage.Bid{
			ExtraPrice: 17,
		},
		storage.Bid{
			ExtraPrice: 15,
		},
	}

	bids = auctionmachine.OrderByDescExtraPrice(bids)
	if bids[0].ExtraPrice != 17 {
		t.Fatalf("Found %v", bids[0].ExtraPrice)
	}
	if bids[2].ExtraPrice != 3 {
		t.Fatalf("Found %v", bids[2].ExtraPrice)
	}
}
