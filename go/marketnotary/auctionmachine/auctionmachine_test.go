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
	machine, err := auctionmachine.NewAuctionMachine(auction, bids, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
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
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() - 10),
		State:      storage.AUCTION_STARTED,
	}
	bids := []storage.Bid{}
	machine, err := auctionmachine.NewAuctionMachine(auction, bids, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
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

func TestStartedAuctionWithBids(t *testing.T) {
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() + 100),
		State:      storage.AUCTION_STARTED,
	}
	bids := []storage.Bid{
		storage.Bid{
			Auction: auction.UUID,
		},
	}
	machine, err := auctionmachine.NewAuctionMachine(auction, bids, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_ASSET_FROZEN {
		t.Fatalf("Expected %v but %v", storage.AUCTION_ASSET_FROZEN, machine.Auction.State)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
}

func TestFrozenAuction(t *testing.T) {
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() + 100),
		State:      storage.AUCTION_ASSET_FROZEN,
	}
	bids := []storage.Bid{
		storage.Bid{
			Auction: auction.UUID,
		},
	}
	machine, err := auctionmachine.NewAuctionMachine(auction, bids, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_ASSET_FROZEN {
		t.Fatalf("Expected %v but %v", storage.AUCTION_ASSET_FROZEN, machine.Auction.State)
	}
}

func TestOutdatedFrozenAuction(t *testing.T) {
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() - 100),
		State:      storage.AUCTION_ASSET_FROZEN,
	}
	bids := []storage.Bid{
		storage.Bid{
			Auction: auction.UUID,
		},
	}
	machine, err := auctionmachine.NewAuctionMachine(auction, bids, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_PAYING {
		t.Fatalf("Expected %v but %v", storage.AUCTION_PAYING, machine.Auction.State)
	}
}
