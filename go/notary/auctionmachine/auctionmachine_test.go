package auctionmachine_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/helper"
	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"

	"github.com/ethereum/go-ethereum/crypto"
	"gotest.tools/assert"
)

func TestAuctionStarted(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		auction := storage.NewAuction()
		auction.ValidUntil = time.Now().Unix() + 100
		auction.PlayerID = "274877906944"
		auction.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e"
		auction.Signature = "381bf58829e11790830eab9924b123d1dbe96dd37b10112729d9d32d476c8d5762598042bb5d5fd63f668455aa3a2ce4e2632241865c26ababa231ad212b5f151b"
		offer := storage.NewOffer()

		m, err := auctionmachine.New(*auction, nil, offer, *bc.Contracts, bc.Owner)
		assert.NilError(t, err)
		assert.NilError(t, m.Process(nil))
		assert.Equal(t, m.StateExtra(), "")
		assert.Equal(t, m.State(), storage.AuctionStarted)
	})

	t.Run("expired", func(t *testing.T) {
		auction := storage.NewAuction()
		auction.ValidUntil = time.Now().Unix() - 10
		offer := storage.NewOffer()

		m, err := auctionmachine.New(*auction, nil, offer, *bc.Contracts, bc.Owner)
		assert.NilError(t, err)
		assert.NilError(t, m.Process(nil))
		assert.Equal(t, m.StateExtra(), "expired")
		assert.Equal(t, m.State(), storage.AuctionEnded)
	})

	t.Run("seller is not the owner", func(t *testing.T) {
		auction := storage.NewAuction()
		auction.ValidUntil = time.Now().Unix() + 100
		auction.PlayerID = "274877906944"
		offer := storage.NewOffer()

		m, err := auctionmachine.New(*auction, nil, offer, *bc.Contracts, bc.Owner)
		assert.NilError(t, err)
		assert.NilError(t, m.Process(nil))
		assert.Equal(t, m.StateExtra(), "seller  is not the owner 0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
		assert.Equal(t, m.State(), storage.AuctionFailed)
	})
}

func TestAuctionStartedGoFrozen(t *testing.T) {
	auction := storage.NewAuction()
	auction.ID = "f1d4501c5158a9018b1618ec4d471c66b663d8f6bffb6e70a0c6584f5c1ea94a"
	auction.ValidUntil = time.Now().Unix() + 100
	auction.PlayerID = "274877906944"
	auction.CurrencyID = 1
	auction.Price = 41234
	auction.Rnd = 4232
	auction.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e"

	playerId, _ := new(big.Int).SetString(auction.PlayerID, 10)
	assert.Assert(t, playerId != nil)

	sellerDigest, err := auction.ComputePutAssetForSaleDigest()

	assert.NilError(t, err)
	signature, err := signer.Sign(sellerDigest.Bytes(), bc.Owner)
	assert.NilError(t, err)
	auction.Signature = hex.EncodeToString(signature)

	// check the seller is the signer
	signer, err := helper.AddressFromHashAndSignature(sellerDigest, signature)

	assert.NilError(t, err)
	assert.Equal(t, signer.Hex(), auction.Seller)

	// check the player owner is the seller
	assert.Assert(t, playerId != nil)
	owner, err := bc.Contracts.Market.GetOwnerPlayer(&bind.CallOpts{}, playerId)
	assert.NilError(t, err)
	assert.Equal(t, owner.Hex(), auction.Seller)

	bid := storage.NewBid()
	bids := []storage.Bid{*bid}
	offer := storage.NewOffer()

	m, err := auctionmachine.New(*auction, bids, offer, *bc.Contracts, bc.Owner)
	assert.NilError(t, err)
	assert.NilError(t, m.Process(nil))
	assert.Equal(t, m.State(), storage.AuctionAssetFrozen)
}

func TestAuctionMachineAllWorkflow(t *testing.T) {
	bc, err := testutils.NewBlockchain()
	assert.NilError(t, err)

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
	validUntil := now + 1000
	offerValidUntil := validUntil + 100
	playerID := big.NewInt(274877906944)
	currencyID := uint8(1)
	price := big.NewInt(41234)
	auctionRnd := big.NewInt(42321)
	extraPrice := big.NewInt(332)
	bidRnd := big.NewInt(1243523)
	teamID := big.NewInt(274877906945)

	hashAuctionMsg, err := signer.ComputePutAssetForSaleDigest(
		currencyID,
		price,
		auctionRnd,
		validUntil,
		offerValidUntil,
		playerID,
	)
	if err != nil {
		t.Fatal(err)
	}
	signAuctionMsg, err := signer.Sign(hashAuctionMsg.Bytes(), alice)
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		ID:         "TheTestingAuction",
		PlayerID:   playerID.String(),
		CurrencyID: int(currencyID),
		Price:      price.Int64(),
		Rnd:        auctionRnd.Int64(),
		ValidUntil: validUntil,
		Signature:  hex.EncodeToString(signAuctionMsg),
		State:      storage.AuctionStarted,
		Seller:     "0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7",
	}

	hashBidMsg, err := signer.HashBidMessage(
		bc.Contracts.Market,
		currencyID,
		price,
		auctionRnd,
		validUntil,
		offerValidUntil,
		playerID,
		extraPrice,
		bidRnd,
		teamID,
	)
	if err != nil {
		t.Fatal(err)
	}
	signBidMsg, err := signer.Sign(hashBidMsg.Bytes(), bob)
	if err != nil {
		t.Fatal(err)
	}
	bids := []storage.Bid{
		storage.Bid{
			AuctionID:  auction.ID,
			ExtraPrice: extraPrice.Int64(),
			Rnd:        bidRnd.Int64(),
			TeamID:     teamID.String(),
			Signature:  hex.EncodeToString(signBidMsg),
			State:      storage.BidAccepted,
		},
	}

	market := v1.NewMockMarketPay()

	machine, err := auctionmachine.New(auction, bids, nil, *bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, machine.State(), storage.AuctionStarted)

	// machine freeze the asset cause of existent bid
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionAssetFrozen)

	// machine do nothing cause auction deadline is not passed
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionAssetFrozen)

	// waiting that the auction deadline
	time.Sleep(10 * time.Second)

	// machine put the auction in paying wait
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionPaying)
	assert.Equal(t, machine.Bids()[0].State, storage.BidAccepted)

	// machine put the bid is paying state
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionPaying)
	assert.Equal(t, machine.Bids()[0].State, storage.BidPaying)
	assert.Assert(t, machine.Bids()[0].PaymentDeadline != 0)

	// set the market transaction as paid
	market.SetOrderStatus(v1.PUBLISHED)

	// machine set the bid to paid and set the auction as withdrable by seller
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.Bids()[0].State, storage.BidPaid)
	assert.Equal(t, machine.State(), storage.AuctionWithdrableBySeller)
	assert.Equal(t, machine.StateExtra(), "")

	t.Run("AuctionWithdrableBySeller", func(t *testing.T) {
		auction := machine.Auction()
		assert.Assert(t, auction.PaymentURL != "")
	})
}

func TestAuctionMachineAllWorkflowWithOffer(t *testing.T) {
	bc, err := testutils.NewBlockchain()
	assert.NilError(t, err)
	tz := uint8(1)
	countryIdx := big.NewInt(0)

	alice, _ := crypto.HexToECDSA("3B878F7892FBBFA30C8bce1DF317C19B853685E707C2CF0EE1927DC516060A54")
	bob, _ := crypto.HexToECDSA("3693a221b147b8888490aa65a86dbef946eccaff76cc1fc93265468822dfb882")
	tx, err := bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		tz,
		countryIdx,
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
		tz,
		countryIdx,
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
	offerValidUntil := now + 100
	validUntil := offerValidUntil + 1000
	playerID := big.NewInt(274877906945)
	currencyID := uint8(1)
	price := big.NewInt(4129834)
	extraPrice := big.NewInt(0)
	dummyRnd := big.NewInt(0)
	offerRnd := big.NewInt(124352439)
	buyerTeamID := big.NewInt(274877906945)

	hashOffer, err := signer.HashBidMessage(
		bc.Contracts.Market,
		currencyID,
		price,
		offerRnd,
		validUntil,
		offerValidUntil,
		playerID,
		extraPrice,
		dummyRnd,
		buyerTeamID,
	)

	signOfferMsg, err := signer.Sign(hashOffer.Bytes(), bob)
	if err != nil {
		t.Fatal(err)
	}

	offer := storage.Offer{
		ID:          "12345",
		PlayerID:    playerID.String(),
		CurrencyID:  int(currencyID),
		Price:       price.Int64(),
		Rnd:         offerRnd.Int64(),
		ValidUntil:  offerValidUntil,
		Signature:   hex.EncodeToString(signOfferMsg),
		State:       storage.OfferStarted,
		StateExtra:  "",
		Seller:      crypto.PubkeyToAddress(alice.PublicKey).Hex(),
		Buyer:       crypto.PubkeyToAddress(bob.PublicKey).Hex(),
		AuctionID:   "",
		BuyerTeamID: buyerTeamID.String(),
	}

	hashAuctionMsg, err := signer.ComputePutAssetForSaleDigest(
		currencyID,
		price,
		offerRnd,
		validUntil,
		offerValidUntil,
		playerID,
	)
	if err != nil {
		t.Fatal(err)
	}
	signAuctionMsg, err := signer.Sign(hashAuctionMsg.Bytes(), alice)
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		ID:         "TheTestingAuction2",
		PlayerID:   playerID.String(),
		CurrencyID: int(currencyID),
		Price:      price.Int64(),
		Rnd:        offerRnd.Int64(),
		ValidUntil: offerValidUntil,
		Signature:  hex.EncodeToString(signAuctionMsg),
		State:      storage.AuctionStarted,
		Seller:     "0x916a407D8cB5B4E533672C908757737F27fE3C25",
	}

	offer.AuctionID = auction.ID
	offer.State = storage.OfferAccepted

	bids := []storage.Bid{
		storage.Bid{
			AuctionID:  auction.ID,
			ExtraPrice: extraPrice.Int64(),
			Rnd:        dummyRnd.Int64(),
			TeamID:     buyerTeamID.String(),
			Signature:  hex.EncodeToString(signOfferMsg),
			State:      storage.BidAccepted,
		},
	}

	market := v1.NewMockMarketPay()

	machine, err := auctionmachine.New(auction, bids, &offer, *bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, machine.State(), storage.AuctionStarted)

	// machine freeze the asset cause of existent bid
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionAssetFrozen)

	// machine do nothing cause auction deadline is not passed
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionAssetFrozen)

	// waiting that the auction deadline
	time.Sleep(10 * time.Second)

	// machine put the auction in paying wait
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionPaying)
	assert.Equal(t, machine.Bids()[0].State, storage.BidAccepted)

	// machine put the bid is paying state
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionPaying)
	assert.Equal(t, machine.Bids()[0].State, storage.BidPaying)
	assert.Assert(t, machine.Bids()[0].PaymentDeadline != 0)

	// set the market transaction as paid
	market.SetOrderStatus(v1.PUBLISHED)

	// machine set the bid to paid and set the auction as withdrable by seller
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.Bids()[0].State, storage.BidPaid)
	assert.Equal(t, machine.State(), storage.AuctionWithdrableBySeller)
	assert.Equal(t, machine.StateExtra(), "")

	t.Run("AuctionWithdrableBySeller", func(t *testing.T) {
		auction := machine.Auction()
		assert.Assert(t, auction.PaymentURL != "")
	})
}
