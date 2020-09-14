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

func TestAcceptOfferReturnTheSignature(t *testing.T) {
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	// We will here assign the next available team to offerer so she can make an offer for the players of a different team
	// we choose that player as a player very far from the current amount of teams (x2)
	nHumanTeams, _ := bc.Contracts.Assets.GetNHumansInCountry(&bind.CallOpts{}, timezoneIdx, countryIdx)
	offererTeamIdx := nHumanTeams.Int64()
	sellerTeamIdx := offererTeamIdx + 1
	// offererTeamId, _ := bc.Contracts.Assets.EncodeTZCountryAndVal(&bind.CallOpts{}, timezoneIdx, countryIdx, big.NewInt(offererTeamIdx))
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

	in := input.AcceptOfferInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+1000, 10)
	in.OfferValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	in.PlayerId = playerId.String()
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	in.OfferId = "12abc345cd"

	hash, err := in.SellerDigest()
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), seller)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	id, err := r.AcceptOffer(struct{ Input input.AcceptOfferInput }{in})
	assert.NilError(t, err)
	id2, err := in.AuctionID()
	assert.NilError(t, err)
	assert.Equal(t, id, id2)
}
