package processor_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/processor"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
	"github.com/google/uuid"
)

func TestOutdatedAuction(t *testing.T) {
	now := time.Now().Unix()
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(now - 10),
		State:      storage.AUCTION_STARTED,
	}
	bids := []storage.Bid{}
	machine := processor.NewAuctionMachine(auction, bids)
	machine.Process()
	if machine.Auction.State != storage.AUCTION_NO_BIDS {
		t.Fatalf("Expected %v but %v", storage.AUCTION_NO_BIDS, auction.State)
	}
}
