package consumer_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/graph-gophers/graphql-go"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

func TestAcceptOffer(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	// The two main time parameters:
	// - the offer should not expire before the seller accepts the offer
	// - then, an auction starts that may end days layer (e.g. 2 days)
	offerValidUntil := big.NewInt(999999999999)
	auctionValidUntil := big.NewInt(999999999999 + 3600*24*2)
	playerId := big.NewInt(274877906940)
	ofererTeamId := big.NewInt(456678987944)

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = offerValidUntil.String()
	inOffer.PlayerId = playerId.String()
	inOffer.TeamId = ofererTeamId.String()
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = 4232
	inOffer.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f"

	extraPrice := big.NewInt(0)
	dummyRnd := big.NewInt(0)

	hashOffer, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		offerValidUntil.Int64(),
		playerId,
		extraPrice,
		dummyRnd,
		ofererTeamId,
		true,
	)
	assert.Equal(t, hashOffer.Hex(), "0x1563f70ce76787ea99b420ad637df3757b492c98cd5a774d7111c861453c270b")
	assert.NilError(t, err)
	offerPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	offerSignature, err := signer.Sign(hashOffer.Bytes(), offerPrivateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(offerSignature), "dbd05f0df6b470d071462ba49956eb472031de84509409823502decb119f2fb36cfb57d5d6f6de5f819731745a4f5533c1805065eebf1a7d56dc9bdced406b231c")
	inOffer.Signature = hex.EncodeToString(offerSignature)

	assert.NilError(t, consumer.CreateOffer(service, tx, inOffer, *bc.Contracts))

	offer, err := service.OfferByRndPrice(tx, inOffer.Rnd, inOffer.Price)
	assert.NilError(t, err)

	in := input.AcceptOfferInput{}
	in.ValidUntil = auctionValidUntil.String()
	in.PlayerId = inOffer.PlayerId
	in.CurrencyId = inOffer.CurrencyId
	in.Price = inOffer.Price
	in.Rnd = inOffer.Rnd
	in.OfferId = graphql.ID(strconv.FormatInt(offer.ID, 10))

	assert.NilError(t, err)
	hash, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		auctionValidUntil.Int64(),
		playerId,
	)
	assert.Equal(t, hash.Hex(), "0x1059367f2fb81d2697a1a2c8bc59e561188e3bb2545bfd783ff58e97f7ec70e4")
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "a063ae70f54381e09eeb4e46f50e52066a4c255945b37a0f0155e541afbc92df7aae2ec4ffc730091013dd062dfb005d255ccf7e39644f1d7d1ac57b214d1cea1c")
	in.Signature = hex.EncodeToString(signature)

	assert.NilError(t, consumer.AcceptOffer(service, tx, in))
	auctionId, err := in.AuctionID()
	assert.NilError(t, err)

	auction, err := service.Auction(tx, string(auctionId))
	assert.NilError(t, err)
	assert.Assert(t, auction != nil)
	assert.Equal(t, auction.Seller, "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
	assert.Equal(t, auction.Price, int64(41234))

	offer, err = service.Offer(tx, offer.ID)
	assert.NilError(t, err)
	assert.Equal(t, string(offer.State), "accepted")
	assert.Equal(t, offer.AuctionID, string(auctionId))

	bids, err := service.Bids(tx, auction.ID)
	assert.Equal(t, string(bids[0].State), "accepted")
	assert.Equal(t, bids[0].Rnd, int64(inOffer.Rnd))
	assert.Equal(t, bids[0].ExtraPrice, int64(0))
	assert.Equal(t, string(bids[0].Signature), hex.EncodeToString(offerSignature))

}

func TestAcceptOfferWithExpiredOffer(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	// in this example the offer expired 5 min ago
	offerValidUntil := big.NewInt(time.Now().Unix() - 5*60)
	auctionValidUntil := big.NewInt(time.Now().Unix() + 48*3600)
	playerId := big.NewInt(274877906940)
	ofererTeamId := big.NewInt(456678987944)

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = offerValidUntil.String()
	inOffer.PlayerId = playerId.String()
	inOffer.TeamId = ofererTeamId.String()
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = 4232
	inOffer.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f"
	assert.NilError(t, err)
	extraPrice := big.NewInt(0)
	dummyRnd := big.NewInt(0)

	hashOffer, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		offerValidUntil.Int64(),
		playerId,
		extraPrice,
		dummyRnd,
		ofererTeamId,
		true,
	)
	// We cannot compare hashes because validUntil is different in every test (it uses time.now())
	assert.NilError(t, err)
	offerPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	offerSignature, err := signer.Sign(hashOffer.Bytes(), offerPrivateKey)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(offerSignature)

	assert.NilError(t, consumer.CreateOffer(service, tx, inOffer, *bc.Contracts))

	offer, err := service.OfferByRndPrice(tx, inOffer.Rnd, inOffer.Price)
	assert.NilError(t, err)

	in := input.AcceptOfferInput{}
	in.ValidUntil = auctionValidUntil.String()
	in.PlayerId = inOffer.PlayerId
	in.CurrencyId = inOffer.CurrencyId
	in.Price = inOffer.Price
	in.Rnd = inOffer.Rnd
	in.OfferId = graphql.ID(strconv.FormatInt(offer.ID, 10))

	assert.NilError(t, err)
	hash, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		auctionValidUntil.Int64(),
		playerId,
	)
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	err = consumer.AcceptOffer(service, tx, in)
	assert.Error(t, err, "Associated Offer is expired")

}

func TestAcceptOfferWithNonExpiredOffer(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	// in this example the offer expires in 5 min
	offerValidUntil := big.NewInt(time.Now().Unix() + 5*60)
	auctionValidUntil := big.NewInt(time.Now().Unix() + 48*3600)
	playerId := big.NewInt(274877906940)
	ofererTeamId := big.NewInt(456678987944)

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = offerValidUntil.String()
	inOffer.PlayerId = playerId.String()
	inOffer.TeamId = ofererTeamId.String()
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = 4232
	inOffer.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f"
	assert.NilError(t, err)
	extraPrice := big.NewInt(0)
	dummyRnd := big.NewInt(0)

	hashOffer, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		offerValidUntil.Int64(),
		playerId,
		extraPrice,
		dummyRnd,
		ofererTeamId,
		true,
	)
	// We cannot compare hashes because validUntil is different in every test (it uses time.now())
	assert.NilError(t, err)
	offerPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	offerSignature, err := signer.Sign(hashOffer.Bytes(), offerPrivateKey)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(offerSignature)

	assert.NilError(t, consumer.CreateOffer(service, tx, inOffer, *bc.Contracts))

	offer, err := service.OfferByRndPrice(tx, inOffer.Rnd, inOffer.Price)
	assert.NilError(t, err)

	in := input.AcceptOfferInput{}
	in.ValidUntil = auctionValidUntil.String()
	in.PlayerId = inOffer.PlayerId
	in.CurrencyId = inOffer.CurrencyId
	in.Price = inOffer.Price
	in.Rnd = inOffer.Rnd
	in.OfferId = graphql.ID(strconv.FormatInt(offer.ID, 10))

	assert.NilError(t, err)
	hash, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		auctionValidUntil.Int64(),
		playerId,
	)
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	err = consumer.AcceptOffer(service, tx, in)
	assert.NilError(t, err)
}
