package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/google/uuid"
)

func TestGetbids(t *testing.T) {
	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	auctionUuid := uuid.New()
	result, err := sto.GetBidsOfAuction(auctionUuid)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 got %v", len(result))
	}
	err = sto.CreateAuction(storage.Auction{
		UUID:       auctionUuid,
		PlayerID:   big.NewInt(5),
		CurrencyID: 1,
		Price:      big.NewInt(3),
		Rnd:        big.NewInt(7),
		ValidUntil: big.NewInt(8),
		Signature:  "0x0",
		State:      storage.AUCTION_STARTED,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = sto.CreateBid(storage.Bid{
		Auction: auctionUuid,
		TeamID:  big.NewInt(2),
		State:   storage.BIDACCEPTED,
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err = sto.GetBidsOfAuction(auctionUuid)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("Expected 1 got %v", len(result))
	}
	bid := result[0]
	if bid.Is2StartAuction != false {
		t.Fatal("Expected false but true")
	}
}

func TestUpdateBidState(t *testing.T) {
	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	auctionUuid := uuid.New()
	result, err := sto.GetBidsOfAuction(auctionUuid)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 got %v", len(result))
	}
	err = sto.CreateAuction(storage.Auction{
		UUID:       auctionUuid,
		PlayerID:   big.NewInt(5),
		CurrencyID: 1,
		Price:      big.NewInt(3),
		Rnd:        big.NewInt(7),
		ValidUntil: big.NewInt(8),
		Signature:  "0x0",
		State:      storage.AUCTION_STARTED,
	})
	if err != nil {
		t.Fatal(err)
	}
	bid := storage.Bid{
		Auction: auctionUuid,
		TeamID:  big.NewInt(2),
		State:   storage.BIDACCEPTED,
	}
	err = sto.CreateBid(bid)
	if err != nil {
		t.Fatal(err)
	}
	bidState := storage. BIDFAILED
	bidStateExtra := "it's just a game dude!"
	err = sto.UpdateBidState(bid.Auction, bid.ExtraPrice, bidState, bidStateExtra)
	if err != nil {
		t.Fatal(err)
	}
	bids, err := sto.GetBidsOfAuction(auctionUuid)
	if err != nil {
		t.Fatal(err)
	}
	if bids[0].State != bidState {
		t.Fatalf("Wrong bid state %v", bids[0].State)
	}
	if bids[0].StateExtra != bidStateExtra {
		t.Fatalf("Wrong bid state extra %v", bidStateExtra)
	}
}

func TestUpdatePaymentId(t *testing.T) {
	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	auctionUuid := uuid.New()
	result, err := sto.GetBidsOfAuction(auctionUuid)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("Expected 0 got %v", len(result))
	}
	err = sto.CreateAuction(storage.Auction{
		UUID:       auctionUuid,
		PlayerID:   big.NewInt(5),
		CurrencyID: 1,
		Price:      big.NewInt(3),
		Rnd:        big.NewInt(7),
		ValidUntil: big.NewInt(8),
		Signature:  "0x0",
		State:      storage.AUCTION_STARTED,
	})
	if err != nil {
		t.Fatal(err)
	}
	bid := storage.Bid{
		Auction: auctionUuid,
		TeamID:  big.NewInt(2),
		State:   storage.BIDACCEPTED,
	}
	err = sto.CreateBid(bid)
	if err != nil {
		t.Fatal(err)
	}
	paymentID := "35565645"
	err = sto.UpdateBidPaymentID(bid.Auction, bid.ExtraPrice, paymentID)
	if err != nil {
		t.Fatal(err)
	}
	bids, err := sto.GetBidsOfAuction(auctionUuid)
	if err != nil {
		t.Fatal(err)
	}
	if bids[0].PaymentID != paymentID {
		t.Fatalf("Wrong paymentID %v", bids[0].PaymentID)
	}
}

// 	result, err = sto.GetBids()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(result) != 1 {
// 		t.Fatalf("Expected 1 got %v", len(result))
// 	}

// 	err = sto.DeleteBet(big.NewInt(1))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err = sto.GetBids()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(result) != 0 {
// 		t.Fatalf("Expected 0 got %v", len(result))
// 	}
// }
