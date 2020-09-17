package auctionmachine_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/helper"
	v1 "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"

	"github.com/ethereum/go-ethereum/crypto"
	"gotest.tools/assert"
)

func TestAuctionStarted(t *testing.T) {
	t.Run("all ok: nothing expired", func(t *testing.T) {
		auction := storage.NewAuction()
		auction.ValidUntil = time.Now().Unix() + 1000
		auction.OfferValidUntil = time.Now().Unix() + 100
		auction.PlayerID = "274877906944" // id of a player in team0 of this timezone (as declared in setup.go)
		auction.Seller = crypto.PubkeyToAddress(bc.Owner.PublicKey).Hex()
		auction.Signature = "381bf58829e11790830eab9924b123d1dbe96dd37b10112729d9d32d476c8d5762598042bb5d5fd63f668455aa3a2ce4e2632241865c26ababa231ad212b5f151b"
		offer := storage.NewOffer()

		m, err := auctionmachine.New(*auction, nil, offer, *bc.Contracts, bc.Owner)
		assert.NilError(t, err)
		assert.NilError(t, m.Process(nil))
		assert.Equal(t, m.StateExtra(), "")
		assert.Equal(t, m.State(), storage.AuctionStarted)
	})

	t.Run("expired validUntil", func(t *testing.T) {
		auction := storage.NewAuction()
		auction.ValidUntil = time.Now().Unix() - 10
		auction.OfferValidUntil = time.Now().Unix() + 100
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
		auction.Seller = "0x03A909262608c650BD9b0ae06E29D90D0F67aC5e"
		offer := storage.NewOffer()

		m, err := auctionmachine.New(*auction, nil, offer, *bc.Contracts, bc.Owner)
		assert.NilError(t, err)
		assert.NilError(t, m.Process(nil))
		assert.Equal(t, m.StateExtra(), "seller 0x03A909262608c650BD9b0ae06E29D90D0F67aC5e is not the owner 0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
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

func TestAuctionMachineAllWorkflow1(t *testing.T) {
	// We will here assign the next available team to offerer so she can make an offer for the players of a different team
	// we choose that player as a player very far from the current amount of teams (x2)
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	offererTeamIdx := nHumanTeams.Int64()
	sellerTeamIdx := offererTeamIdx + 1
	offererTeamId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, big.NewInt(offererTeamIdx))
	offerer, _ := crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	seller, _ := crypto.HexToECDSA("0A878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, (big.NewInt(2 + 18*sellerTeamIdx)))

	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(offerer.PublicKey),
	)
	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(seller.PublicKey),
	)

	now := time.Now().Unix()
	validUntil := now + 8
	offerValidUntil := int64(0)
	playerID := playerId
	currencyID := uint8(1)
	price := big.NewInt(41234)
	auctionRnd := big.NewInt(42321)
	extraPrice := big.NewInt(332)
	bidRnd := big.NewInt(1243523)
	teamID := offererTeamId

	sellerDigest, err := signer.ComputePutAssetForSaleDigest(
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
	signAuctionMsg, err := signer.Sign(sellerDigest.Bytes(), seller)
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		PlayerID:        playerID.String(),
		CurrencyID:      int(currencyID),
		Price:           price.Int64(),
		Rnd:             auctionRnd.Int64(),
		ValidUntil:      validUntil,
		OfferValidUntil: offerValidUntil,
		Signature:       hex.EncodeToString(signAuctionMsg),
		State:           storage.AuctionStarted,
		Seller:          crypto.PubkeyToAddress(seller.PublicKey).Hex(),
	}
	auction.ID, _ = auction.ComputeID()

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
	signBidMsg, err := signer.Sign(hashBidMsg.Bytes(), offerer)
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
	// We will here assign the next available team to offerer so she can make an offer for the players of a different team
	// we choose that player as a player very far from the current amount of teams (x2)
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	offererTeamIdx := nHumanTeams.Int64()
	sellerTeamIdx := offererTeamIdx + 1
	offererTeamId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, big.NewInt(offererTeamIdx))
	offerer, _ := crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	seller, _ := crypto.HexToECDSA("0A878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, (big.NewInt(2 + 18*sellerTeamIdx)))

	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(offerer.PublicKey),
	)
	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(seller.PublicKey),
	)

	now := time.Now().Unix()
	offerValidUntil := now + 5
	validUntil := offerValidUntil + 11
	playerID := playerId.String()
	currencyID := uint8(1)
	price := big.NewInt(4129834)
	extraPrice := big.NewInt(0)
	dummyRnd := big.NewInt(0)
	offerRnd := big.NewInt(124352439)
	buyerTeamID := offererTeamId

	// validUntilStr := strconv.FormatInt(validUntil, 10)
	offerValidUntilStr := strconv.FormatInt(offerValidUntil, 10)

	inOffer := input.CreateOfferInput{}
	// inOffer.Signature   string
	inOffer.PlayerId = playerID
	inOffer.CurrencyId = int32(currencyID)
	inOffer.Price = int32(price.Int64())
	inOffer.Rnd = int32(offerRnd.Int64())
	inOffer.ValidUntil = offerValidUntilStr
	inOffer.BuyerTeamId = offererTeamId.String()
	inOffer.Seller = crypto.PubkeyToAddress(seller.PublicKey).Hex()

	hashOffer, err := inOffer.Hash(*bc.Contracts)

	signOfferMsg, err := signer.Sign(hashOffer.Bytes(), offerer)
	if err != nil {
		t.Fatal(err)
	}

	offerId, _ := inOffer.ID(*bc.Contracts)

	offer := storage.Offer{
		PlayerID:    inOffer.PlayerId,
		CurrencyID:  int(currencyID),
		Price:       price.Int64(),
		Rnd:         offerRnd.Int64(),
		ValidUntil:  offerValidUntil,
		Signature:   hex.EncodeToString(signOfferMsg),
		State:       storage.OfferStarted,
		StateExtra:  "",
		Seller:      crypto.PubkeyToAddress(seller.PublicKey).Hex(),
		Buyer:       crypto.PubkeyToAddress(offerer.PublicKey).Hex(),
		AuctionID:   string(offerId),
		BuyerTeamID: buyerTeamID.String(),
	}

	hashAuctionMsg, err := signer.ComputePutAssetForSaleDigest(
		currencyID,
		price,
		offerRnd,
		validUntil,
		offerValidUntil,
		playerId,
	)
	if err != nil {
		t.Fatal(err)
	}
	signAuctionMsg, err := signer.Sign(hashAuctionMsg.Bytes(), seller)
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		ID:              offer.AuctionID,
		PlayerID:        offer.PlayerID,
		CurrencyID:      offer.CurrencyID,
		Price:           offer.Price,
		Rnd:             offer.Rnd,
		ValidUntil:      validUntil,
		OfferValidUntil: offer.ValidUntil,
		Signature:       hex.EncodeToString(signAuctionMsg),
		State:           storage.AuctionStarted,
		Seller:          crypto.PubkeyToAddress(seller.PublicKey).Hex(),
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
	assert.Equal(t, machine.Bids()[0].State, storage.BidAccepted)

	// machine freeze the asset cause of existent bid
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionAssetFrozen)
	assert.Equal(t, machine.Bids()[0].State, storage.BidAccepted)

	// machine do nothing cause auction deadline is not passed
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionAssetFrozen)
	assert.Equal(t, machine.Bids()[0].State, storage.BidAccepted)

	// waiting that the auction deadline
	time.Sleep(20 * time.Second)

	// machine put the auction in paying wait
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionPaying)
	assert.Equal(t, machine.Bids()[0].State, storage.BidAccepted)

	// machine put the bid in paying state
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
