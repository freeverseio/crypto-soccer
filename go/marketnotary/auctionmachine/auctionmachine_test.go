package auctionmachine_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/freeverseio/crypto-soccer/go/marketnotary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
	"github.com/google/uuid"
)

func TestAuctionWithNoBids(t *testing.T) {
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() + 100),
		State:      storage.AUCTION_STARTED,
	}
	bids := []storage.Bid{}
	machine := auctionmachine.NewAuctionMachine(auction, bids, nil, nil)
	err := machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_STARTED {
		t.Fatalf("Expected %v but %v", storage.AUCTION_STARTED, machine.Auction.State)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuctionOutdatedWithNoBids(t *testing.T) {
	// bc, err := testutils.NewBlockchainNode()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// err = bc.DeployContracts(bc.Owner)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// timezoneIdx := uint8(1)
	// err = bc.InitOneTimezone(timezoneIdx)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() - 10),
		State:      storage.AUCTION_STARTED,
	}
	bids := []storage.Bid{}
	machine := auctionmachine.NewAuctionMachine(auction, bids, nil, nil)
	err := machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_NO_BIDS {
		t.Fatalf("Expected %v but %v", storage.AUCTION_NO_BIDS, machine.Auction.State)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
}
