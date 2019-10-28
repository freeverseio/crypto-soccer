package auctionmachine_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
	"github.com/google/uuid"
)

func TestStartedStateToNoBids(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}

	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	timezoneIdx := uint8(1)
	err = bc.InitOneTimezone(timezoneIdx)
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() - 10),
		State:      storage.AUCTION_STARTED,
	}
	bids := []storage.Bid{}
	machine := auctionmachine.NewAuctionMachine(auction, bids, bc.Market, nil)
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_NO_BIDS {
		t.Fatalf("Expected %v but %v", storage.AUCTION_NO_BIDS, auction.State)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
}
