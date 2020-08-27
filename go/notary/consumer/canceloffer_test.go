package consumer_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"github.com/graph-gophers/graphql-go"
	"gotest.tools/assert"
)

func TestCancelOffer(t *testing.T) {
	service := postgres.NewStorageService(db)
	assert.NilError(t, service.Begin())
	defer service.Rollback()

	in := input.CreateOfferInput{}
	in.ValidUntil = "999999999999"
	in.PlayerId = "274877906940"
	in.BuyerTeamId = "456678987944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 4232
	in.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f"
	playerId, _ := new(big.Int).SetString(in.PlayerId, 10)
	teamId, _ := new(big.Int).SetString(in.BuyerTeamId, 10)
	validUntil, err := strconv.ParseInt(in.ValidUntil, 10, 64)
	dummyRnd := big.NewInt(0)
	offerExtraPrice := big.NewInt(0)
	isOffer := true
	assert.NilError(t, err)
	// an offer cannot be signed with non null extraPrice:
	hash, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		validUntil,
		playerId,
		big.NewInt(2),
		dummyRnd,
		teamId,
		isOffer,
	)
	assert.Error(t, err, "offers must have zero extraPrice")
	// an offer cannot be signed with non null bid.Rnd:
	hash, err = signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		validUntil,
		playerId,
		offerExtraPrice,
		big.NewInt(2),
		teamId,
		isOffer,
	)
	assert.Error(t, err, "offers must have zero bidRnd")
	// it should now work:
	hash, err = signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		validUntil,
		playerId,
		offerExtraPrice,
		dummyRnd,
		teamId,
		isOffer,
	)
	assert.NilError(t, err)
	assert.Equal(t, hash.Hex(), "0x1563f70ce76787ea99b420ad637df3757b492c98cd5a774d7111c861453c270b")
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "dbd05f0df6b470d071462ba49956eb472031de84509409823502decb119f2fb36cfb57d5d6f6de5f819731745a4f5533c1805065eebf1a7d56dc9bdced406b231c")
	in.Signature = hex.EncodeToString(signature)

	assert.NilError(t, consumer.CreateOffer(service, in, *bc.Contracts))
	assert.NilError(t, err)

	offer, err := service.OfferByRndPrice(in.Rnd, in.Price)
	assert.NilError(t, err)
	assert.Assert(t, offer != nil)
	assert.Equal(t, offer.Seller, "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f")

	// cancel the offer
	inCancel := input.CancelOfferInput{}
	inCancel.OfferId = graphql.ID(offer.ID)
	err = consumer.CancelOffer(service, inCancel)
	assert.NilError(t, err)

	offer, err = service.OfferByRndPrice(in.Rnd, in.Price)
	assert.Equal(t, string(offer.State), "cancelled")
}
