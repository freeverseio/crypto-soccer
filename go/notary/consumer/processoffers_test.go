package consumer_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

func TestProcessOffers(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	in := input.CreateOfferInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()-100000, 10)
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
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	assert.NilError(t, consumer.CreateOffer(service, tx, in, *bc.Contracts))
	assert.NilError(t, err)

	offer, err := service.OfferByRndPrice(tx, in.Rnd, in.Price)
	assert.NilError(t, err)
	assert.Assert(t, offer != nil)
	assert.Equal(t, offer.Seller, "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f")

	// Process pending offers
	offers, err := service.OfferPendingOffers(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(offers), 1)

	err = consumer.ProcessOffers(service, tx)
	assert.NilError(t, err)

	offers, err = service.OfferPendingOffers(tx)
	assert.NilError(t, err)
	assert.Equal(t, len(offers), 0)
}
