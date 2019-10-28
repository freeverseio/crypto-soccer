package auctionmachine_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
	"github.com/google/uuid"
)

func TestAuctionWithNoBids(t *testing.T) {
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() + 100),
		State:      storage.AUCTION_STARTED,
	}
	bids := []storage.Bid{}
	machine, err := auctionmachine.New(auction, bids, nil, nil)
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
	machine, err := auctionmachine.New(auction, bids, nil, nil)
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
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	err = bc.InitOneTimezone(1)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := bc.Leagues.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		1,
		big.NewInt(0),
		common.HexToAddress("0x9c0511bfc917d6395E3ccEe7F15A1898FC5CCd97"),
	)
	err = bc.WaitReceipt(tx, 5)
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   big.NewInt(274877906944),
		CurrencyID: uint8(1),
		Price:      big.NewInt(41234),
		Rnd:        big.NewInt(42321),
		ValidUntil: big.NewInt(1572284694),
		Signature:  "0x025bbed3e0810e682ad500d9f35c90246e7580bbc44ccc81aec951636d2b7dd228c27239aa2fb7ef4e1b729f89f9ccf1152897949f22b9f35f30706c1f39f4791b",
		State:      storage.AUCTION_STARTED,
	}
	auctionHiddenPrice, err := bc.Market.HashPrivateMsg(&bind.CallOpts{}, auction.CurrencyID, auction.Price, auction.Rnd)
	if err != nil {
		t.Fatal(err)
	}
	result := hex.EncodeToString(auctionHiddenPrice[:])
	if result != "4200de738160a9e6b8f69648fbb7feb323f73fac5acff1b7bb546bb7ac3591fa" {
		t.Fatalf("Expected 4200de738160a9e6b8f69648fbb7feb323f73fac5acff1b7bb546bb7ac3591fa got %v", result)
	}

	bids := []storage.Bid{
		storage.Bid{
			Auction: auction.UUID,
		},
	}
	machine, err := auctionmachine.New(auction, bids, bc.Market, bc.Owner)
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
	machine, err := auctionmachine.New(auction, bids, nil, nil)
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
	machine, err := auctionmachine.New(auction, bids, nil, nil)
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

func TestPayingAuction(t *testing.T) {
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() - 1),
		State:      storage.AUCTION_PAYING,
	}
	bids := []storage.Bid{
		storage.Bid{
			Auction: auction.UUID,
		},
	}
	machine, err := auctionmachine.New(auction, bids, nil, nil)
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

func TestPayingPaymentDoneAuction(t *testing.T) {
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() - 70),
		State:      storage.AUCTION_PAYING,
	}
	bids := []storage.Bid{
		storage.Bid{
			Auction: auction.UUID,
		},
	}
	machine, err := auctionmachine.New(auction, bids, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State == storage.AUCTION_PAYING {
		t.Fatalf("Expected not %v", machine.Auction.State)
	}
}
