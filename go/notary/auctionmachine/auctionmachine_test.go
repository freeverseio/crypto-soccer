package auctionmachine_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	"github.com/freeverseio/crypto-soccer/go/helper"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	"gotest.tools/assert"
)

func newTestMarket() *marketpay.MarketPay {
	market, err := marketpay.New()
	if err != nil {
		panic(err)
	}
	return market
}

func TestAuctionStarted(t *testing.T) {
	market := newMarket(t)

	auction := storage.NewAuction()
	auction.ValidUntil = time.Now().Unix() + 100
	m, err := auctionmachine.New(*auction, nil, bc.Contracts, bc.Owner)
	assert.NilError(t, err)
	assert.Equal(t, m.State(), storage.AUCTION_STARTED)
	assert.NilError(t, err)
	assert.NilError(t, m.Process(market))
	assert.Equal(t, m.State(), storage.AUCTION_STARTED)

	auction.ValidUntil = time.Now().Unix() - 10
	m, err = auctionmachine.New(*auction, nil, bc.Contracts, bc.Owner)
	assert.NilError(t, err)
	assert.Equal(t, m.State(), storage.AUCTION_STARTED)
	assert.NilError(t, err)
	assert.NilError(t, m.Process(market))
	assert.Equal(t, m.State(), storage.AUCTION_NO_BIDS)
}

func TestAuctionWithNoBids(t *testing.T) {
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: time.Now().Unix() + 100,
		State:      storage.AUCTION_STARTED,
	}
	bids := []*storage.Bid{}
	machine, err := auctionmachine.New(auction, bids, bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	market := newTestMarket()
	err = machine.Process(market)
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_STARTED {
		t.Fatalf("Expected %v but %v", storage.AUCTION_STARTED, machine.Auction.State)
	}
	err = machine.Process(market)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuctionOutdatedWithNoBids(t *testing.T) {
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: time.Now().Unix() - 10,
		State:      storage.AUCTION_STARTED,
	}
	bids := []*storage.Bid{}
	machine, err := auctionmachine.New(auction, bids, bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	market := newTestMarket()
	err = machine.Process(market)
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_NO_BIDS {
		t.Fatalf("Expected %v but %v", storage.AUCTION_NO_BIDS, machine.Auction.State)
	}
	err = machine.Process(market)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartedAuctionWithBids(t *testing.T) {
	tx, err := bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		1,
		big.NewInt(0),
		common.HexToAddress("0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7"),
	)
	_, err = helper.WaitReceipt(bc.Client, tx, 5)
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   big.NewInt(274877906944),
		CurrencyID: uint8(1),
		Price:      big.NewInt(41234),
		Rnd:        big.NewInt(42321),
		ValidUntil: 2000000000,
		Signature:  "0x4cc92984c7ee4fe678b0c9b1da26b6757d9000964d514bdaddc73493393ab299276bad78fd41091f9fe6c169adaa3e8e7db146a83e0a2e1b60480320443919471c",
		State:      storage.AUCTION_STARTED,
	}
	bids := []*storage.Bid{
		&storage.Bid{
			Auction: auction.UUID,
		},
	}
	machine, err := auctionmachine.New(auction, bids, bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	market := newTestMarket()
	err = machine.Process(market)
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_FAILED {
		t.Fatalf("Expected %v but %v", storage.AUCTION_FAILED, machine.Auction.State)
	}
}

func TestFrozenAuction(t *testing.T) {
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: time.Now().Unix() + 100,
		State:      storage.AUCTION_ASSET_FROZEN,
	}
	bids := []*storage.Bid{
		&storage.Bid{
			Auction: auction.UUID,
		},
	}
	machine, err := auctionmachine.New(auction, bids, bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	market := newTestMarket()
	err = machine.Process(market)
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
		ValidUntil: time.Now().Unix() - 100,
		State:      storage.AUCTION_ASSET_FROZEN,
	}
	bids := []*storage.Bid{
		&storage.Bid{
			Auction: auction.UUID,
		},
	}
	machine, err := auctionmachine.New(auction, bids, bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	market := newTestMarket()
	err = machine.Process(market)
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
		ValidUntil: time.Now().Unix() - 1,
		State:      storage.AUCTION_PAYING,
	}
	bids := []*storage.Bid{
		&storage.Bid{
			Auction: auction.UUID,
			State:   storage.BIDACCEPTED,
		},
	}
	machine, err := auctionmachine.New(auction, bids, bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	market := newTestMarket()
	err = machine.Process(market)
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_PAYING {
		t.Fatalf("Expected %v but %v", storage.AUCTION_PAYING, machine.Auction.State)
	}
}

func TestPayingPaymentDoneAuction(t *testing.T) {
	alice, _ := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	bob, _ := crypto.HexToECDSA("3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882")

	tx, err := bc.Contracts.Assets.TransferFirstBotToAddr(
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
	tx, err = bc.Contracts.Assets.TransferFirstBotToAddr(
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
	validUntil := now + 8
	playerID := big.NewInt(274877906944)
	currencyID := uint8(1)
	price := big.NewInt(41234)
	auctionRnd := big.NewInt(42321)
	extraPrice := big.NewInt(332)
	bidRnd := big.NewInt(1243523)
	teamID := big.NewInt(274877906945)
	isOffer2StartAuction := false

	signer := signer.NewSigner(bc.Contracts, nil)
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
	auction := storage.Auction{
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
	machine, err := auctionmachine.New(auction, bids, bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	market := newTestMarket()
	err = machine.Process(market)
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_ASSET_FROZEN {
		t.Fatalf("Expected not %v", machine.Auction.State)
	}
	time.Sleep(10 * time.Second)
	err = machine.Process(market)
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_PAYING {
		t.Fatalf("Expected not %v", machine.Auction.State)
	}
	if machine.Bids[0].State != storage.BIDACCEPTED {
		t.Fatalf("Expected not %v", machine.Bids[0].State)
	}
	err = machine.Process(market)
	if err != nil {
		t.Fatal(err)
	}
	if machine.Auction.State != storage.AUCTION_PAYING {
		t.Fatalf("Expected not %v", machine.Auction.State)
	}
	if machine.Bids[0].State != storage.BIDPAYING {
		t.Fatalf("Expected not %v", machine.Bids[0].State)
	}
	if machine.Bids[0].PaymentDeadline == 0 {
		t.Fatalf("Wrong bid timeout %v", machine.Bids[0].PaymentDeadline)
	}
	// following is commented because we need an action from the user to make marketpay set it as PAID
	// time.Sleep(10 * time.Second)
	// err = machine.Process(market)
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
