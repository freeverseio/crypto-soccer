package bidmachine_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/bidmachine"
	"github.com/freeverseio/crypto-soccer/go/testutils"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func TestNotPayingAuction(t *testing.T) {
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	auction := &storage.Auction{State: storage.AUCTION_ASSET_FROZEN}
	bid := &storage.Bid{}
	_, err = bidmachine.New(
		auction,
		bid,
		bc.Contracts,
		bc.Owner,
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
	auction := &storage.Auction{State: storage.AUCTION_PAYING}
	bid := &storage.Bid{}
	_, err = bidmachine.New(
		auction,
		bid,
		bc.Contracts,
		bc.Owner,
	)
	if err != nil {
		t.Fatalf("Not accepting %v auction", auction.State)
	}
}

func TestFirstAlive(t *testing.T) {
	bid := bidmachine.FirstAlive(nil)
	if bid != nil {
		t.Fatalf("Wrong result: %v", bid)
	}
	bids := []*storage.Bid{}
	bid = bidmachine.FirstAlive(bids)
	if bid != nil {
		t.Fatalf("Wrong result: %v", bid)
	}
	bids = []*storage.Bid{&storage.Bid{State: storage.BIDFAILED}}
	bid = bidmachine.FirstAlive(bids)
	if bid != nil {
		t.Fatalf("Wrong result: %v", bid)
	}
	bids = append(bids, &storage.Bid{State: storage.BIDACCEPTED, ExtraPrice: 10})
	bid = bidmachine.FirstAlive(bids)
	if bid != bids[1] {
		t.Fatalf("Expected %v result: %v", bids[0], bid)
	}
	bids = append(bids, &storage.Bid{State: storage.BIDACCEPTED, ExtraPrice: 11})
	bid = bidmachine.FirstAlive(bids)
	if bid != bids[2] {
		t.Fatalf("Wrong result: %v", bid)
	}
	bids = append(bids, &storage.Bid{State: storage.BIDPAYING, ExtraPrice: 11})
	bid = bidmachine.FirstAlive(bids)
	if bid != bids[3] {
		t.Fatalf("Wrong result: %v", bid)
	}
}

func TestExpiredBidNoTransit(t *testing.T) {
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	auction := &storage.Auction{State: storage.AUCTION_PAYING}
	bid := &storage.Bid{State: storage.BIDFAILED}
	machine, err := bidmachine.New(
		auction,
		bid,
		bc.Contracts,
		bc.Owner,
	)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if bid.State != storage.BIDFAILED {
		t.Fatalf("Wrong state %v", bid.State)
	}
}

func TestAcceptBidTransitToPaying(t *testing.T) {
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	auction := &storage.Auction{
		State:      storage.AUCTION_PAYING,
		ValidUntil: big.NewInt(10),
	}
	bid := &storage.Bid{
		State:           storage.BIDACCEPTED,
		PaymentDeadline: big.NewInt(3),
	}
	machine, err := bidmachine.New(
		auction,
		bid,
		bc.Contracts,
		bc.Owner,
	)
	if err != nil {
		t.Fatal(err)
	}
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if bid.State != storage.BIDPAYING {
		t.Fatalf("Wrong state %v", bid.State)
	}
	if bid.PaymentDeadline.String() != "21610" {
		t.Fatalf("Wrong deadline %v", bid.PaymentDeadline)
	}
}
