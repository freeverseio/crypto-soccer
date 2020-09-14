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
	"github.com/graph-gophers/graphql-go"
	"gotest.tools/assert"
)

func TestCreateOffer1(t *testing.T) {
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

	ch := make(chan interface{}, 10)

	// We use offerValidUntil for offers, and validUntil for accept offer and make bids later
	validUntil := time.Now().Unix() + 1000
	offerValidUntil := time.Now().Unix() + 100

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = strconv.FormatInt(offerValidUntil, 10)
	inOffer.PlayerId = playerId.String()
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = int32(42321)
	inOffer.BuyerTeamId = offererTeamId.String()
	inOffer.Seller = crypto.PubkeyToAddress(seller.PublicKey).Hex()
	offerID, _ := inOffer.ID(*bc.Contracts)

	hash, err := inOffer.Hash(*bc.Contracts)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), offerer)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(signature)

	mockOffer := storage.Offer{
		ID:          string(offerID),
		PlayerID:    inOffer.PlayerId,
		CurrencyID:  int(inOffer.CurrencyId),
		Price:       int64(inOffer.Price),
		Rnd:         int64(inOffer.Rnd),
		ValidUntil:  validUntil,
		Signature:   "0x" + inOffer.Signature,
		State:       storage.OfferStarted,
		StateExtra:  "",
		Seller:      inOffer.Seller,
		Buyer:       crypto.PubkeyToAddress(offerer.PublicKey).Hex(),
		AuctionID:   "",
		BuyerTeamID: inOffer.BuyerTeamId,
	}
	mockOffersByPlayerId := []storage.Offer{mockOffer}

	mock := mockup.Tx{
		AuctionInsertFunc:      func(auction storage.Auction) error { return nil },
		AuctionsByPlayerIdFunc: func(ID string) ([]storage.Auction, error) { return []storage.Auction{}, nil },
		OfferInsertFunc:        func(offer storage.Offer) error { return nil },
		BidInsertFunc:          func(bid storage.Bid) error { return nil },
		CommitFunc:             func() error { return nil },
		OffersByPlayerIdFunc:   func(playerId string) ([]storage.Offer, error) { return mockOffersByPlayerId, nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)

	_, err = r.CreateOffer(struct{ Input input.CreateOfferInput }{inOffer})

	// When you accept the offer, validUntil is redefined, and offerValidUntil is inherited from the offer
	acceptOfferIn := input.AcceptOfferInput{}
	acceptOfferIn.OfferValidUntil = inOffer.ValidUntil
	acceptOfferIn.ValidUntil = strconv.FormatInt(validUntil, 10)
	acceptOfferIn.PlayerId = inOffer.PlayerId
	acceptOfferIn.CurrencyId = inOffer.CurrencyId
	acceptOfferIn.Price = inOffer.Price
	acceptOfferIn.Rnd = inOffer.Rnd
	acceptOfferIn.OfferId = graphql.ID(string(offerID))

	sellerDigest, err := acceptOfferIn.SellerDigest()
	signature, err = signer.Sign(sellerDigest.Bytes(), seller)
	assert.NilError(t, err)
	acceptOfferIn.Signature = hex.EncodeToString(signature)

	_, err = r.AcceptOffer(struct{ Input input.AcceptOfferInput }{acceptOfferIn})
	assert.NilError(t, err)

	// The original offer signature should be valid to create an auction
	auctionId, err := acceptOfferIn.AuctionID()
	inBid := input.CreateBidInput{}
	inBid.AuctionId = auctionId
	inBid.ExtraPrice = 0
	inBid.Rnd = 0
	inBid.TeamId = inOffer.BuyerTeamId
	inBid.Signature = inOffer.Signature

	err = r.CreateBid(struct{ Input input.CreateBidInput }{inBid})
	assert.NilError(t, err)
}

func TestCreateOfferSignedByNotOwnedPlayer(t *testing.T) {
	// identical to previous test but choosing a playerId that belongs to offerer's team too
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)

	// We will here assign the next available team to offerer so she can make an offer for the players of a different team
	// we choose that player as a player very far from the current amount of teams (x2)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	offererTeamIdx := nHumanTeams.Int64()
	sellerTeamIdx := offererTeamIdx + 1
	offererTeamId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, big.NewInt(offererTeamIdx))
	offerer, _ := crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	seller, _ := crypto.HexToECDSA("0A878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, (big.NewInt(-2 + 18*sellerTeamIdx)))

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

	ch := make(chan interface{}, 10)

	// We use offerValidUntil for offers, and validUntil for accept offer and make bids later
	validUntil := time.Now().Unix() + 1000
	offerValidUntil := time.Now().Unix() + 100

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = strconv.FormatInt(offerValidUntil, 10)
	inOffer.PlayerId = playerId.String()
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = int32(42321)
	inOffer.BuyerTeamId = offererTeamId.String()
	offerID, _ := inOffer.ID(*bc.Contracts)

	hash, err := inOffer.Hash(*bc.Contracts)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), offerer)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(signature)

	mockOffer := storage.Offer{
		ID:          string(offerID),
		PlayerID:    inOffer.PlayerId,
		CurrencyID:  int(inOffer.CurrencyId),
		Price:       int64(inOffer.Price),
		Rnd:         int64(inOffer.Rnd),
		ValidUntil:  validUntil,
		Signature:   "0x" + inOffer.Signature,
		State:       storage.OfferStarted,
		StateExtra:  "",
		Seller:      inOffer.Seller,
		Buyer:       crypto.PubkeyToAddress(offerer.PublicKey).Hex(),
		AuctionID:   "",
		BuyerTeamID: inOffer.BuyerTeamId,
	}
	mockOffersByPlayerId := []storage.Offer{mockOffer}

	mock := mockup.Tx{
		AuctionInsertFunc:      func(auction storage.Auction) error { return nil },
		AuctionsByPlayerIdFunc: func(ID string) ([]storage.Auction, error) { return []storage.Auction{}, nil },
		OfferInsertFunc:        func(offer storage.Offer) error { return nil },
		BidInsertFunc:          func(bid storage.Bid) error { return nil },
		CommitFunc:             func() error { return nil },
		OffersByPlayerIdFunc:   func(playerId string) ([]storage.Offer, error) { return mockOffersByPlayerId, nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CreateOffer(struct{ Input input.CreateOfferInput }{inOffer})

	// When you accept the offer, validUntil is redefined, and offerValidUntil is inherited from the offer
	acceptOfferIn := input.AcceptOfferInput{}
	acceptOfferIn.OfferValidUntil = inOffer.ValidUntil
	acceptOfferIn.ValidUntil = strconv.FormatInt(validUntil, 10)
	acceptOfferIn.PlayerId = inOffer.PlayerId
	acceptOfferIn.CurrencyId = inOffer.CurrencyId
	acceptOfferIn.Price = inOffer.Price
	acceptOfferIn.Rnd = inOffer.Rnd
	acceptOfferIn.OfferId = graphql.ID(string(offerID))

	sellerDigest, err := acceptOfferIn.SellerDigest()
	signature, err = signer.Sign(sellerDigest.Bytes(), seller)
	assert.NilError(t, err)
	acceptOfferIn.Signature = hex.EncodeToString(signature)

	_, err = r.AcceptOffer(struct{ Input input.AcceptOfferInput }{acceptOfferIn})
	errString := err.Error()
	l := intMin(len(errString)-1, 35)
	assert.Equal(t, errString[:l], "signer is not the owner of playerId")
}

func intMin(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestCreateOfferSameOwner(t *testing.T) {
	// offerer tries to offer and sell his own player
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)

	// We will here assign the next available team to offerer so she can make an offer for the players of a different team
	// we choose that player as a player very far from the current amount of teams (x2)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	offererTeamIdx := nHumanTeams.Int64()
	// sellerTeamIdx := offererTeamIdx + 1
	offererTeamId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, big.NewInt(offererTeamIdx))
	offerer, _ := crypto.HexToECDSA("9B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	seller, _ := crypto.HexToECDSA("0A878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	playerId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, (big.NewInt(2 + 18*offererTeamIdx)))

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

	ch := make(chan interface{}, 10)

	// We use offerValidUntil for offers, and validUntil for accept offer and make bids later
	validUntil := time.Now().Unix() + 1000
	offerValidUntil := time.Now().Unix() + 100

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = strconv.FormatInt(offerValidUntil, 10)
	inOffer.PlayerId = playerId.String()
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = int32(42321)
	inOffer.BuyerTeamId = offererTeamId.String()
	offerID, _ := inOffer.ID(*bc.Contracts)

	hash, err := inOffer.Hash(*bc.Contracts)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), offerer)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(signature)

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

	// When you accept the offer, validUntil is redefined, and offerValidUntil is inherited from the offer
	acceptOfferIn := input.AcceptOfferInput{}
	acceptOfferIn.BuyerTeamId = inOffer.BuyerTeamId
	acceptOfferIn.OfferValidUntil = inOffer.ValidUntil
	acceptOfferIn.ValidUntil = strconv.FormatInt(validUntil, 10)
	acceptOfferIn.PlayerId = inOffer.PlayerId
	acceptOfferIn.CurrencyId = inOffer.CurrencyId
	acceptOfferIn.Price = inOffer.Price
	acceptOfferIn.Rnd = inOffer.Rnd
	acceptOfferIn.OfferId = graphql.ID(string(offerID))

	sellerDigest, err := acceptOfferIn.SellerDigest()
	signature, err = signer.Sign(sellerDigest.Bytes(), offerer)
	assert.NilError(t, err)
	acceptOfferIn.Signature = hex.EncodeToString(signature)

	_, err = r.AcceptOffer(struct{ Input input.AcceptOfferInput }{acceptOfferIn})

	assert.Error(t, err, "the buyerTeam already owns the player it is making an offer for")
}

func TestCreateOfferMadeByNotTeamOwner(t *testing.T) {
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)

	// We will here assign the next available team to offerer so she can make an offer for the players of a different team
	// we choose that player as a player very far from the current amount of teams (x2)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	offererTeamIdx := nHumanTeams.Int64()
	sellerTeamIdx := offererTeamIdx + 1
	offererTeamId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, big.NewInt(offererTeamIdx+3))
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

	ch := make(chan interface{}, 10)

	// We use offerValidUntil for offers, and validUntil for accept offer and make bids later
	validUntil := time.Now().Unix() + 1000
	offerValidUntil := time.Now().Unix() + 100

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = strconv.FormatInt(offerValidUntil, 10)
	inOffer.PlayerId = playerId.String()
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = int32(42321)
	inOffer.BuyerTeamId = offererTeamId.String()
	offerID, _ := inOffer.ID(*bc.Contracts)

	hash, err := inOffer.Hash(*bc.Contracts)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), offerer)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(signature)

	mockOffer := storage.Offer{
		ID:          string(offerID),
		PlayerID:    inOffer.PlayerId,
		CurrencyID:  int(inOffer.CurrencyId),
		Price:       int64(inOffer.Price),
		Rnd:         int64(inOffer.Rnd),
		ValidUntil:  validUntil,
		Signature:   "0x" + inOffer.Signature,
		State:       storage.OfferStarted,
		StateExtra:  "",
		Seller:      inOffer.Seller,
		Buyer:       crypto.PubkeyToAddress(offerer.PublicKey).Hex(),
		AuctionID:   "",
		BuyerTeamID: inOffer.BuyerTeamId,
	}
	mockOffersByPlayerId := []storage.Offer{mockOffer}

	mock := mockup.Tx{
		AuctionInsertFunc:      func(auction storage.Auction) error { return nil },
		AuctionsByPlayerIdFunc: func(ID string) ([]storage.Auction, error) { return []storage.Auction{}, nil },
		OfferInsertFunc:        func(offer storage.Offer) error { return nil },
		BidInsertFunc:          func(bid storage.Bid) error { return nil },
		CommitFunc:             func() error { return nil },
		OffersByPlayerIdFunc:   func(playerId string) ([]storage.Offer, error) { return mockOffersByPlayerId, nil },
	}
	service := &mockup.StorageService{
		BeginFunc: func() (storage.Tx, error) { return &mock, nil },
	}

	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, service)
	_, err = r.CreateOffer(struct{ Input input.CreateOfferInput }{inOffer})

	// When you accept the offer, validUntil is redefined, and offerValidUntil is inherited from the offer
	acceptOfferIn := input.AcceptOfferInput{}
	acceptOfferIn.OfferValidUntil = inOffer.ValidUntil
	acceptOfferIn.ValidUntil = strconv.FormatInt(validUntil, 10)
	acceptOfferIn.PlayerId = inOffer.PlayerId
	acceptOfferIn.CurrencyId = inOffer.CurrencyId
	acceptOfferIn.Price = inOffer.Price
	acceptOfferIn.Rnd = inOffer.Rnd
	acceptOfferIn.OfferId = graphql.ID(string(offerID))

	sellerDigest, err := acceptOfferIn.SellerDigest()
	signature, err = signer.Sign(sellerDigest.Bytes(), seller)
	assert.NilError(t, err)
	acceptOfferIn.Signature = hex.EncodeToString(signature)

	_, err = r.AcceptOffer(struct{ Input input.AcceptOfferInput }{acceptOfferIn})

	assert.NilError(t, err)

	// The original offer signature should be valid to create an auction
	auctionId, err := acceptOfferIn.AuctionID()
	inBid := input.CreateBidInput{}
	inBid.AuctionId = auctionId
	inBid.ExtraPrice = 0
	inBid.Rnd = 0
	inBid.TeamId = inOffer.BuyerTeamId
	inBid.Signature = inOffer.Signature

	err = r.CreateBid(struct{ Input input.CreateBidInput }{inBid})
	errString := err.Error()
	l := intMin(len(errString)-1, 33)
	assert.Equal(t, errString[:l], "signer is not the owner of teamId")
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
	assert.Equal(t, hash.Hex(), "0x0bd108f7e6cf8ae86312d1680aececc07dd8c4d97f67879ca159b3ba5a074a90")
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "5f04ac4d856b9a33b9590e23589b711b69f0a2f8d26bcb47bbf9e9a48d8761c116e449c3e9307c672e0b0811a0dde30a33a067dacd9a8cfa83838263344a78a21c")
	inOffer.Signature = hex.EncodeToString(signature)

	_, err = r.CreateOffer(struct{ Input input.CreateOfferInput }{inOffer})
	assert.Error(t, err, "signer is not the owner of teamId 456678987944")
}
