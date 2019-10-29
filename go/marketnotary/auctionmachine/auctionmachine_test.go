package auctionmachine_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/marketnotary/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
	"github.com/google/uuid"
)

func TestAuctionWithNoBids(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() + 100),
		State:      storage.AUCTION_STARTED,
	}
	bids := []storage.Bid{}
	machine, err := auctionmachine.New(auction, bids, bc.Market, bc.Owner, bc.Client)
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
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() - 10),
		State:      storage.AUCTION_STARTED,
	}
	bids := []storage.Bid{}
	machine, err := auctionmachine.New(auction, bids, bc.Market, bc.Owner, bc.Client)
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
		common.HexToAddress("0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7"),
	)
	err = helper.WaitReceipt(bc.Client, tx, 5)
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   big.NewInt(274877906944),
		CurrencyID: uint8(1),
		Price:      big.NewInt(41234),
		Rnd:        big.NewInt(42321),
		ValidUntil: big.NewInt(2000000000),
		Signature:  "0x4cc92984c7ee4fe678b0c9b1da26b6757d9000964d514bdaddc73493393ab299276bad78fd41091f9fe6c169adaa3e8e7db146a83e0a2e1b60480320443919471c",
		State:      storage.AUCTION_STARTED,
	}
	bids := []storage.Bid{
		storage.Bid{
			Auction: auction.UUID,
		},
	}
	machine, err := auctionmachine.New(auction, bids, bc.Market, bc.Owner, bc.Client)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_FAILED_TO_FREEZE {
		t.Fatalf("Expected %v but %v", storage.AUCTION_FAILED_TO_FREEZE, machine.Auction.State)
	}
}

func TestFrozenAuction(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
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
	machine, err := auctionmachine.New(auction, bids, bc.Market, bc.Owner, bc.Client)
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
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
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
	machine, err := auctionmachine.New(auction, bids, bc.Market, bc.Owner, bc.Client)
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
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
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
	machine, err := auctionmachine.New(auction, bids, bc.Market, bc.Owner, bc.Client)
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
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
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
	machine, err := auctionmachine.New(auction, bids, bc.Market, bc.Owner, bc.Client)
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
