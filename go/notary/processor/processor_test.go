package processor_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/freeverseio/crypto-soccer/go/notary/processor"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
)

func TestOutdatedAuction(t *testing.T) {
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().Unix()
	auction := storage.Auction{
		UUID:       uuid.New(),
		ValidUntil: big.NewInt(now - 10),
		State:      storage.AUCTION_STARTED,
	}
	err = sto.CreateAuction(auction)
	if err != nil {
		t.Fatal(err)
	}
	processor, err := processor.NewProcessor(sto, bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	err = processor.Process()
	if err != nil {
		t.Fatal(err)
	}
	auctions, err := sto.GetAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if auctions[0].State != storage.AUCTION_NO_BIDS {
		t.Fatalf("Expected %v but %v", storage.AUCTION_NO_BIDS, auctions[0].State)
	}
}

func TestAuctionWithBid(t *testing.T) {
	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	bc, err := testutils.NewBlockchainNode()
	if err != nil {
		t.Fatal(err)
	}
	err = bc.DeployContracts(bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().Unix()
	auction := storage.Auction{
		UUID:       uuid.New(),
		PlayerID:   big.NewInt(65),
		Price:      big.NewInt(7),
		Rnd:        big.NewInt(73),
		ValidUntil: big.NewInt(now + 100),
		Signature:  "0x4cc92984c7ee4fe678b0c9b1da26b6757d9000964d514bdaddc73493393ab299276bad78fd41091f9fe6c169adaa3e8e7db146a83e0a2e1b60480320443919471c",
		State:      storage.AUCTION_STARTED,
	}
	err = sto.CreateAuction(auction)
	if err != nil {
		t.Fatal(err)
	}
	bid := storage.Bid{
		Auction:         auction.UUID,
		State:           storage.BIDACCEPTED,
		PaymentDeadline: big.NewInt(0),
	}
	err = sto.CreateBid(bid)
	if err != nil {
		t.Fatal(err)
	}
	processor, err := processor.NewProcessor(sto, bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}
	err = processor.Process()
	if err != nil {
		t.Fatal(err)
	}
	auctions, err := sto.GetAuctions()
	if err != nil {
		t.Fatal(err)
	}
	if auctions[0].State != storage.AUCTION_FAILED {
		t.Fatalf("Expected %v but %v", storage.AUCTION_FAILED, auctions[0].State)
	}
	bids, err := sto.GetBidsOfAuction(auctions[0].UUID)
	if err != nil {
		t.Fatal(err)
	}
	if bids[0].State != storage.BIDACCEPTED {
		t.Fatalf("Expects %v got %v", storage.BIDACCEPTED, bids[0].State)
	}
}

// func TestFreezePlayer(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	ganache := testutils.NewGanache()
// 	alice, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	value := new(big.Int)
// 	value.SetString("50000000000000000000", 10)
// 	_, err = ganache.TransferWei(value, ganache.Owner, ganache.Public(alice))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	processor, err := processor.NewProcessor(sto, ganache.Client, ganache.Market, ganache.Owner)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	timezone := uint8(1)
// 	countryIdxInTZ := big.NewInt(0)
// 	teamId0, err := ganache.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezone, countryIdxInTZ, big.NewInt(0))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = ganache.Assets.TransferFirstBotToAddr(bind.NewKeyedTransactor(alice), timezone, countryIdxInTZ, crypto.PubkeyToAddress(alice.PublicKey))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	team0PlayerIds, err := ganache.Assets.GetPlayerIdsInTeam(&bind.CallOpts{}, teamId0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	team0PlayerId0 := team0PlayerIds[0]
// 	if team0PlayerId0.String() != "274877906944" {
// 		t.Fatalf("Wrong player id : %v", team0PlayerId0.String())
// 	}

// 	Auction := storage.Auction{
// 		PlayerId:   team0PlayerId0,
// 		CurrencyId: 1,
// 		Price:      big.NewInt(41234),
// 		Rnd:        big.NewInt(42321),
// 		ValidUntil: big.NewInt(2000000000),
// 		Signature:  "0x4cc92984c7ee4fe678b0c9b1da26b6757d9000964d514bdaddc73493393ab299276bad78fd41091f9fe6c169adaa3e8e7db146a83e0a2e1b60480320443919471c",
// 	}
// 	err = processor.FreezePlayer(Auction)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestProcess(t *testing.T) {
// 	sto, err := storage.NewSqlite3("../../../market.db/00_schema.sql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	ganache := testutils.NewGanache()
// 	alice, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
// 	bob, err := crypto.HexToECDSA("3693a221b147b7338490aa65a86dbef946eccaff76cc1fc93265468822dfb882")

// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	value := new(big.Int)
// 	value.SetString("50000000000000000000", 10)
// 	_, err = ganache.TransferWei(value, ganache.Owner, ganache.Public(alice))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = ganache.TransferWei(value, ganache.Owner, ganache.Public(bob))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	processor, err := processor.NewProcessor(sto, ganache.Client, ganache.Market, ganache.Owner)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	timezone := uint8(1)
// 	countryIdxInTZ := big.NewInt(0)
// 	teamId0, err := ganache.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezone, countryIdxInTZ, big.NewInt(0))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = ganache.Assets.TransferFirstBotToAddr(bind.NewKeyedTransactor(alice), timezone, countryIdxInTZ, crypto.PubkeyToAddress(alice.PublicKey))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	teamId1, err := ganache.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezone, countryIdxInTZ, big.NewInt(1))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_, err = ganache.Assets.TransferFirstBotToAddr(bind.NewKeyedTransactor(bob), timezone, countryIdxInTZ, crypto.PubkeyToAddress(bob.PublicKey))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	team0PlayerIds, err := ganache.Assets.GetPlayerIdsInTeam(&bind.CallOpts{}, teamId0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	validUntil := big.NewInt(2000000000)
// 	playerId := team0PlayerIds[0]
// 	typeOfTX := uint8(1)
// 	currencyId := uint8(1)
// 	price := big.NewInt(41234)
// 	rnd := big.NewInt(42321)
// 	teamId := teamId1

// 	originOwner := ganache.GetPlayerOwner(playerId)
// 	if originOwner != ganache.Public(alice) {
// 		t.Fatalf("Expected originOwner ALICE but got %v", originOwner)
// 	}
// 	sto.CreateAuction(storage.Auction{
// 		PlayerId:   playerId,
// 		CurrencyId: currencyId,
// 		Price:      price,
// 		Rnd:        rnd,
// 		ValidUntil: validUntil,
// 		TypeOfTx:   typeOfTX,
// 		Signature:  "0xac466c2139f6edce74d18161252922d8368dce25c3e508de98e8659e9a994a000dd33bd3034aea26fe99b54b1df240041f77afb0a2be508a83e7d35482b20a951c",
// 	})
// 	processor.Process()
// 	targetOwner := ganache.GetPlayerOwner(playerId)
// 	if targetOwner != crypto.PubkeyToAddress(alice.PublicKey) {
// 		t.Fatalf("Expectedf originOwner ALICE but got %v", targetOwner)
// 	}

// 	sto.CreateBet(storage.Bid{
// 		PlayerId:  playerId,
// 		TeamId:    teamId,
// 		Signature: "0x44bb117064e1e2a8ef5fed99f3ec9281f95ef7caea595db2c36071963f74e4c904e8c61d6cb75aaef61718e1d2dff49bc3c55c886e7b3d3e73db31a1af3c61721b",
// 	})

// 	processor.Process()
// 	targetOwner = ganache.GetPlayerOwner(playerId)
// 	if targetOwner != crypto.PubkeyToAddress(bob.PublicKey) {
// 		t.Fatalf("Expected originOwner BOB but got %v", targetOwner)
// 	}

// 	bids, err := sto.GetBids()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(bids) != 0 {
// 		t.Fatalf("Expercted 0 but got %v", len(bids))
// 	}
// 	Auctions, err := sto.GetAuctions()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(Auctions) != 0 {
// 		t.Fatalf("Expercted 0 but got %v", len(Auctions))
// 	}
// }
