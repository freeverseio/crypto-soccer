package gql_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/mockup"
	"gotest.tools/assert"
)

func TestCreateOffer1(t *testing.T) {
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)

	// We will here assign the next available team to offerer so she can make an offer for the players at team = 0
	// playerId from the second team is made an offer
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, big.NewInt(nHumanTeams.Int64()*18+10))
	teamId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, big.NewInt(nHumanTeams.Int64()))
	offerer, _ := crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")

	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(offerer.PublicKey),
	)

	ch := make(chan interface{}, 10)

	mock := mockup.Tx{
		AuctionInsertFunc:      func(auction storage.Auction) error { return nil },
		AuctionsByPlayerIdFunc: func(ID string) ([]storage.Auction, error) { return []storage.Auction{}, nil },
		OfferInsertFunc:        func(offer storage.Offer) error { return nil },
		BidInsertFunc:          func(bid storage.Bid) error { return nil },
		CommitFunc:             func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	// We use offerValidUntil for offers, and validUntil for accept offer and make bids later
	validUntil := time.Now().Unix() + 1000
	offerValidUntil := time.Now().Unix() + 100

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = strconv.FormatInt(offerValidUntil, 10)
	inOffer.PlayerId = playerId.String()
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = int32(42321)
	inOffer.BuyerTeamId = teamId.String()

	hash, err := inOffer.Hash(*bc.Contracts)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), offerer)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(signature)

	// When you accept the offer, validUntil is redefined, and offerValidUntil is inherited from the offer
	acceptOfferIn := input.AcceptOfferInput{}
	acceptOfferIn.OfferValidUntil = inOffer.ValidUntil
	acceptOfferIn.ValidUntil = strconv.FormatInt(validUntil, 10)
	acceptOfferIn.PlayerId = inOffer.PlayerId
	acceptOfferIn.CurrencyId = inOffer.CurrencyId
	acceptOfferIn.Price = inOffer.Price
	acceptOfferIn.Rnd = inOffer.Rnd

	sellerDigest, err := acceptOfferIn.SellerDigest()
	signature, err = signer.Sign(sellerDigest.Bytes(), offerer)
	assert.NilError(t, err)
	acceptOfferIn.Signature = hex.EncodeToString(signature)

	_, err = r.CreateAuctionFromOffer(struct {
		Input input.AcceptOfferInput
	}{acceptOfferIn})
	assert.NilError(t, err)

	// The original offer signature should be valid to create an auction
	auctionId, err := acceptOfferIn.AuctionID()
	inBid := input.CreateBidInput{}
	inBid.AuctionId = auctionId
	inBid.ExtraPrice = 0
	inBid.Rnd = 0
	inBid.TeamId = inOffer.BuyerTeamId
	inBid.Signature = inOffer.Signature

	_, err = r.CreateBid(struct{ Input input.CreateBidInput }{inBid})
	assert.NilError(t, err)
}

func TestCreateOfferSameOwner(t *testing.T) {
	offerer := bc.Owner
	offererRnd := int32(42321)
	offerValidUntil := time.Now().Unix() + 100

	mock := mockup.Tx{
		AuctionInsertFunc:      func(auction storage.Auction) error { return nil },
		AuctionsByPlayerIdFunc: func(ID string) ([]storage.Auction, error) { return []storage.Auction{}, nil },
		CommitFunc:             func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = strconv.FormatInt(offerValidUntil, 10)
	inOffer.PlayerId = "274877906944"
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = offererRnd
	inOffer.BuyerTeamId = "274877906945"
	teamId, _ := new(big.Int).SetString(inOffer.BuyerTeamId, 10)
	playerId, _ := new(big.Int).SetString(inOffer.PlayerId, 10)
	validUntil, err := strconv.ParseInt(inOffer.ValidUntil, 10, 64)

	dummyValidUntil := int64(0)
	dummyExtraPrice := big.NewInt(0)
	dummyRnd := big.NewInt(0)
	hashOffer, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		dummyValidUntil,
		validUntil,
		playerId,
		dummyExtraPrice,
		dummyRnd,
		teamId,
	)

	assert.NilError(t, err)
	signatureOffer, err := signer.Sign(hashOffer.Bytes(), offerer)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(signatureOffer)
	_, err = r.CreateOffer(struct{ Input input.CreateOfferInput }{inOffer})
	assert.Error(t, err, "signer is the owner of playerId 274877906944 you can't make an offer for your player")

}

func TestCreateOfferNotTeamOwner(t *testing.T) {
	teamOwnedByOffered := big.NewInt(274877906945)
	teamNotOwnedByOffered := big.NewInt(274877906948)

	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	offerer, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(offerer.PublicKey),
	)

	offererRnd := int32(42321)
	offerValidUntil := time.Now().Unix() + 100

	ch := make(chan interface{}, 10)
	mock := mockup.Tx{
		AuctionInsertFunc:      func(auction storage.Auction) error { return nil },
		AuctionsByPlayerIdFunc: func(ID string) ([]storage.Auction, error) { return []storage.Auction{}, nil },
		OfferInsertFunc:        func(offer storage.Offer) error { return nil },
		CommitFunc:             func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = strconv.FormatInt(offerValidUntil, 10)
	inOffer.PlayerId = "274877906944"
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = offererRnd
	inOffer.BuyerTeamId = teamNotOwnedByOffered.String()
	playerId, _ := new(big.Int).SetString(inOffer.PlayerId, 10)
	validUntil, err := strconv.ParseInt(inOffer.ValidUntil, 10, 64)

	dummyValidUntil := int64(0)
	dummyExtraPrice := big.NewInt(0)
	dummyRnd := big.NewInt(0)
	hashOffer, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		dummyValidUntil,
		validUntil,
		playerId,
		dummyExtraPrice,
		dummyRnd,
		teamNotOwnedByOffered,
	)

	assert.NilError(t, err)
	signatureOffer, err := signer.Sign(hashOffer.Bytes(), offerer)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(signatureOffer)
	_, err = r.CreateOffer(struct{ Input input.CreateOfferInput }{inOffer})
	assert.Error(t, err, "signer is not the owner of teamId 274877906948")

	// exactly same call but with a team truly owned by offerer
	inOffer.BuyerTeamId = teamOwnedByOffered.String()

	hashOffer, err = signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		dummyValidUntil,
		validUntil,
		playerId,
		dummyExtraPrice,
		dummyRnd,
		teamOwnedByOffered,
	)
	assert.NilError(t, err)
	signatureOffer, err = signer.Sign(hashOffer.Bytes(), offerer)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(signatureOffer)
	_, err = r.CreateOffer(struct{ Input input.CreateOfferInput }{inOffer})
	assert.NilError(t, err)
}

func TestCreateOfferExConsumer(t *testing.T) {
	ch := make(chan interface{}, 10)
	mock := mockup.Tx{
		AuctionInsertFunc:      func(auction storage.Auction) error { return nil },
		AuctionsByPlayerIdFunc: func(ID string) ([]storage.Auction, error) { return []storage.Auction{}, nil },
		OfferInsertFunc:        func(offer storage.Offer) error { return nil },
		CommitFunc:             func() error { return nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = "999999999999"
	inOffer.PlayerId = "274877906940"
	inOffer.BuyerTeamId = "456678987944"
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = 4232
	inOffer.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f"
	playerId, _ := new(big.Int).SetString(inOffer.PlayerId, 10)
	teamId, _ := new(big.Int).SetString(inOffer.BuyerTeamId, 10)
	validUntil, err := strconv.ParseInt(inOffer.ValidUntil, 10, 64)
	dummyRnd := big.NewInt(0)
	offerExtraPrice := big.NewInt(0)
	assert.NilError(t, err)

	dummyValidUntil := int64(0)
	// an offer cannot be signed with non null extraPrice:
	_, err = signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		dummyValidUntil,
		validUntil,
		playerId,
		big.NewInt(2),
		dummyRnd,
		teamId,
	)
	assert.Error(t, err, "offers must have zero extraPrice")

	// an offer cannot be signed with non null bid.Rnd:
	_, err = signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		dummyValidUntil,
		validUntil,
		playerId,
		offerExtraPrice,
		big.NewInt(2),
		teamId,
	)
	assert.Error(t, err, "offers must have zero bidRnd")
	// it should now work:
	hash, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		dummyValidUntil,
		validUntil,
		playerId,
		offerExtraPrice,
		dummyRnd,
		teamId,
	)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x1563f70ce76787ea99b420ad637df3757b492c98cd5a774d7111c861453c270b")
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "dbd05f0df6b470d071462ba49956eb472031de84509409823502decb119f2fb36cfb57d5d6f6de5f819731745a4f5533c1805065eebf1a7d56dc9bdced406b231c")
	inOffer.Signature = hex.EncodeToString(signature)

	_, err = r.CreateOffer(struct{ Input input.CreateOfferInput }{inOffer})
	assert.Error(t, err, "signer is not the owner of teamId 456678987944")
}
