package input_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/helper"
	marketpay "github.com/freeverseio/crypto-soccer/go/marketpay/v1"
	"github.com/freeverseio/crypto-soccer/go/notary/auctionmachine"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/testutils"
	"gotest.tools/assert"
)

func TestCreateOfferInputHash(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	in.BuyerTeamId = "20"
	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x48280aaca3224b385bcc4e4b94662cbf17f989f99426943da0e1a10cd2e5a4a0")
}

func TestCreateOfferSignerAddress(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "10"
	in.CurrencyId = 1
	in.Price = 41234
	in.BuyerTeamId = "20"
	in.Rnd = 42321

	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x48280aaca3224b385bcc4e4b94662cbf17f989f99426943da0e1a10cd2e5a4a0")

	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)

	in.Signature = hex.EncodeToString(signature)
	address, err := in.SignerAddress(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, address.Hex(), crypto.PubkeyToAddress(bc.Owner.PublicKey).Hex())
}

func TestSignerOfOfferIsOwnerOfTeam(t *testing.T) {
	tz := uint8(1)
	countryIdxInTz := big.NewInt(0)

	// first team (teamdIdx = 0) is assigned to during setup to owner (players 0...17),
	// We will here assign the next available team to alice so she can put a playerId for sale
	// playerId from the second team is made an offer
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, tz, countryIdxInTz)
	teamId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, tz, countryIdxInTz, big.NewInt(nHumanTeams.Int64()))
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, tz, countryIdxInTz, big.NewInt(3))
	alice, _ := crypto.HexToECDSA("6B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	// The offer fails because alice is not the owner of the team. We then transfer the team to Alice, and offer works.
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = playerId.String()
	in.CurrencyId = 1
	in.BuyerTeamId = teamId.String()
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), alice)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	isOwner, err := in.IsSignerOwnerOfTeam(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, isOwner, false)

	tx, err := bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		tz,
		countryIdxInTz,
		crypto.PubkeyToAddress(alice.PublicKey),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = helper.WaitReceipt(bc.Client, tx, 5)
	if err != nil {
		t.Fatal(err)
	}

	isOwner, err = in.IsSignerOwnerOfTeam(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, isOwner, true)
}

func TestCreateOfferGetOwner(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.BuyerTeamId = "20"
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	owner, err := in.GetOwner(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, crypto.PubkeyToAddress(bc.Owner.PublicKey).Hex(), owner.Hex())
}

func TestCreateOfferPlayerFrozen(t *testing.T) {

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
	validUntil := now + 800
	offerValidUntil := now + 8
	playerID := big.NewInt(274877906944)
	currencyID := uint8(1)
	price := big.NewInt(41234)
	offererRnd := big.NewInt(42321)
	extraPrice := big.NewInt(332)
	bidRnd := big.NewInt(1243523)
	teamID := big.NewInt(274877906945)

	dummyValidUntil := int64(0)
	dummyExtraPrice := big.NewInt(0)
	dummyRnd := big.NewInt(0)
	// First the offer creates an offer by signing an apropriate bid message
	hashBidMsg, err := signer.HashBidMessage(
		bc.Contracts.Market,
		currencyID,
		price,
		offererRnd,
		dummyValidUntil,
		offerValidUntil,
		playerID,
		dummyExtraPrice,
		dummyRnd,
		teamID,
	)
	if err != nil {
		t.Fatal(err)
	}
	signBidMsg, err := signer.Sign(hashBidMsg.Bytes(), bob)
	if err != nil {
		t.Fatal(err)
	}
	// the auctionId (if the offer is accepted) is determined directly from the offer message
	auctionId, err := signer.ComputeAuctionId(
		currencyID,
		price,
		offererRnd,
		dummyValidUntil,
		offerValidUntil,
		playerID,
	)
	// The seller can accept the offer, and adds a new validUntil
	sellerDigest, err := signer.ComputePutAssetForSaleDigest(
		currencyID,
		price,
		offererRnd,
		validUntil,
		offerValidUntil,
		playerID,
	)
	if err != nil {
		t.Fatal(err)
	}
	signAuctionMsg, err := signer.Sign(sellerDigest.Bytes(), alice)
	if err != nil {
		t.Fatal(err)
	}
	auction := storage.Auction{
		ID:              auctionId.Hex(),
		PlayerID:        playerID.String(),
		CurrencyID:      int(currencyID),
		Price:           price.Int64(),
		Rnd:             offererRnd.Int64(),
		ValidUntil:      validUntil,
		OfferValidUntil: offerValidUntil,
		Signature:       "0x" + hex.EncodeToString(signAuctionMsg),
		State:           storage.AuctionStarted,
		Seller:          "0x291081e5a1bF0b9dF6633e4868C88e1FA48900e7",
	}

	bids := []storage.Bid{
		storage.Bid{
			AuctionID:  auction.ID,
			ExtraPrice: extraPrice.Int64(),
			Rnd:        bidRnd.Int64(),
			TeamID:     teamID.String(),
			Signature:  "0x" + hex.EncodeToString(signBidMsg),
			State:      storage.BidAccepted,
		},
	}

	market := marketpay.NewMockMarketPay()
	offer := storage.NewOffer()

	machine, err := auctionmachine.New(auction, bids, offer, *bc.Contracts, bc.Owner)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, machine.State(), storage.AuctionStarted)

	// machine freeze the asset cause of existent bid
	assert.NilError(t, machine.Process(market))
	assert.Equal(t, machine.State(), storage.AuctionAssetFrozen)

	// try to create offer which will fail because asset is frozen
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.BuyerTeamId = "20"
	in.Price = 41234
	in.Rnd = 42321

	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	isPlayerFrozen, err := in.IsPlayerFrozen(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, true, isPlayerFrozen)
}

// func TestCreateOfferPlayerAlreadyOnSale(t *testing.T) {
// 	service := postgres.NewStorageService(db)
// 	tx, err := service.Begin()
// 	assert.NilError(t, err)
// 	defer tx.Rollback()

// 	in := input.CreateAuctionInput{}
// 	in.ValidUntil = "999999999999"
// 	in.PlayerId = "274877906944"
// 	in.CurrencyId = 1
// 	in.Price = 41234
// 	in.Rnd = 4232
// 	playerId, _ := new(big.Int).SetString(in.PlayerId, 10)
// 	validUntil, err := strconv.ParseInt(in.ValidUntil, 10, 64)
// 	assert.NilError(t, err)
// 	hash, err := signer.HashSellMessage(
// 		uint8(in.CurrencyId),
// 		big.NewInt(int64(in.Price)),
// 		big.NewInt(int64(in.Rnd)),
// 		validUntil,
// 		playerId,
// 	)
// 	assert.Equal(t, hash.Hex(), "0xf1d4501c5158a9018b1618ec4d471c66b663d8f6bffb6e70a0c6584f5c1ea94a")
// 	assert.NilError(t, err)
// 	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
// 	assert.NilError(t, err)
// 	signature, err := signer.Sign(hash.Bytes(), privateKey)
// 	assert.NilError(t, err)
// 	assert.Equal(t, hex.EncodeToString(signature), "381bf58829e11790830eab9924b123d1dbe96dd37b10112729d9d32d476c8d5762598042bb5d5fd63f668455aa3a2ce4e2632241865c26ababa231ad212b5f151b")
// 	in.Signature = hex.EncodeToString(signature)

// 	assert.NilError(t, consumer.CreateAuction(tx, in))

// 	// try to create offer which will fail because asset is on sale
// 	inOffer := input.CreateOfferInput{}
// 	inOffer.ValidUntil = "2000000000"
// 	inOffer.PlayerId = "274877906944"
// 	inOffer.CurrencyId = 1
// 	inOffer.BuyerTeamId = "20"
// 	inOffer.Price = 41234
// 	inOffer.Rnd = 42321

// 	hashOffer, err := inOffer.Hash(*bc.Contracts)
// 	assert.NilError(t, err)
// 	signatureOffer, err := signer.Sign(hashOffer.Bytes(), bc.Owner)
// 	assert.NilError(t, err)
// 	inOffer.Signature = hex.EncodeToString(signatureOffer)
// 	isPlayerOnSale, err := inOffer.IsPlayerOnSale(*bc.Contracts, tx)
// 	assert.NilError(t, err)
// 	assert.Equal(t, true, isPlayerOnSale)
// }

func TestCreateOfferInputHashBigIntPlayer(t *testing.T) {
	in := input.CreateOfferInput{}
	in.ValidUntil = "2000000000"
	in.PlayerId = "25723578238440869144533393071649442553899076447028039543423578"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	in.BuyerTeamId = "20"
	hash, err := in.Hash(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x84347c60199d0444e4ec8bbdf788ab1afff4a59f02c146a8c493d2f73955d7c6")
}
