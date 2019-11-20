package auctionmachine_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/signer"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
	"github.com/google/uuid"
)

func TestAuctionWithNoBids(t *testing.T) {
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	auction := &storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() + 100),
		State:      storage.AUCTION_STARTED,
	}
	bids := []*storage.Bid{}
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
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	auction := &storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() - 10),
		State:      storage.AUCTION_STARTED,
	}
	bids := []*storage.Bid{}
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
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	tx, err := bc.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		1,
		big.NewInt(0),
		common.HexToAddress("0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7"),
	)
	_, err = helper.WaitReceipt(bc.Client, tx, 5)
	if err != nil {
		t.Fatal(err)
	}
	auction := &storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   big.NewInt(274877906944),
		CurrencyID: uint8(1),
		Price:      big.NewInt(41234),
		Rnd:        big.NewInt(42321),
		ValidUntil: big.NewInt(2000000000),
		Signature:  "0x4cc92984c7ee4fe678b0c9b1da26b6757d9000964d514bdaddc73493393ab299276bad78fd41091f9fe6c169adaa3e8e7db146a83e0a2e1b60480320443919471c",
		State:      storage.AUCTION_STARTED,
	}
	bids := []*storage.Bid{
		&storage.Bid{
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
	auction := &storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() + 100),
		State:      storage.AUCTION_ASSET_FROZEN,
	}
	bids := []*storage.Bid{
		&storage.Bid{
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
	auction := &storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() - 100),
		State:      storage.AUCTION_ASSET_FROZEN,
	}
	bids := []*storage.Bid{
		&storage.Bid{
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
	auction := &storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(time.Now().Unix() - 1),
		State:      storage.AUCTION_PAYING,
	}
	bids := []*storage.Bid{
		&storage.Bid{
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
	bc, err := testutils.NewBlockchainNodeDeployAndInit()
	if err != nil {
		t.Fatal(err)
	}
	alice, _ := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	bob, _ := crypto.HexToECDSA("3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882")

	tx, err := bc.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		1,
		big.NewInt(0),
		crypto.PubkeyToAddress(alice.PublicKey),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = helper.WaitReceipt(bc.Client, tx, 5)
	if err != nil {
		t.Fatal(err)
	}
	tx, err = bc.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		1,
		big.NewInt(0),
		crypto.PubkeyToAddress(bob.PublicKey),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = helper.WaitReceipt(bc.Client, tx, 5)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().Unix()
	validUntil := big.NewInt(now + 8)
	playerID := big.NewInt(274877906944)
	currencyID := uint8(1)
	price := big.NewInt(41234)
	auctionRnd := big.NewInt(42321)
	extraPrice := big.NewInt(332)
	bidRnd := big.NewInt(1243523)
	teamID := big.NewInt(274877906945)
	isOffer2StartAuction := false

	signer := signer.NewSigner(bc.Market, nil)
	hashAuctionMsg, err := signer.HashSellMessage(
		currencyID,
		price,
		auctionRnd,
		validUntil,
		playerID,
	)
	if err != nil {
		t.Fatal(err)
	}
	signAuctionMsg, err := signer.Sign(hashAuctionMsg, alice)
	if err != nil {
		t.Fatal(err)
	}
	auction := &storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   playerID,
		CurrencyID: currencyID,
		Price:      price,
		Rnd:        auctionRnd,
		ValidUntil: validUntil,
		Signature:  "0x" + hex.EncodeToString(signAuctionMsg),
		State:      storage.AUCTION_STARTED,
	}

	hashBidMsg, err := signer.HashBidMessage(
		currencyID,
		price,
		auctionRnd,
		validUntil,
		playerID,
		extraPrice,
		bidRnd,
		teamID,
		isOffer2StartAuction,
	)
	if err != nil {
		t.Fatal(err)
	}
	signBidMsg, err := signer.Sign(hashBidMsg, bob)
	if err != nil {
		t.Fatal(err)
	}
	bids := []*storage.Bid{
		&storage.Bid{
			Auction:    auction.UUID,
			ExtraPrice: extraPrice.Int64(),
			Rnd:        bidRnd.Int64(),
			TeamID:     teamID,
			Signature:  "0x" + hex.EncodeToString(signBidMsg),
			State:      storage.BIDACCEPTED,
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
		t.Fatalf("Expected not %v", machine.Auction.State)
	}
	time.Sleep(10 * time.Second)
	err = machine.Process()
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_PAYING {
		t.Fatalf("Expected not %v", machine.Auction.State)
	}
	if machine.Bids[0].State != storage.BIDPAYING {
		t.Fatalf("Expected not %v", machine.Bids[0].State)
	}
	if machine.Auction.State != storage.AUCTION_PAYING {
		t.Fatalf("Expected not %v", machine.Auction.State)
	}
	if machine.Bids[0].State != storage.BIDPAYING {
		t.Fatalf("Expected not %v", machine.Bids[0].State)
	}
	// following is commented because we need an action from the user to make marketpay set it as PAID
	// time.Sleep(10 * time.Second)
	// err = machine.Process()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if machine.Auction.State != storage.AUCTION_PAID {
	// 	t.Fatalf("Expected not %v", machine.Auction.State)
	// }
	// if machine.Bids[0].State != storage.BIDPAID {
	// 	t.Fatalf("Expected not %v", machine.Bids[0].State)
	// }
}
