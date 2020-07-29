package consumer_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
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

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = "999999999999"
	inOffer.PlayerId = "274877906940"
	inOffer.TeamId = "456678987944"
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = 4232
	inOffer.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f"

	extraPrice := big.NewInt(0)
	dummyRnd := big.NewInt(0)
	offerPlayerId, _ := new(big.Int).SetString(inOffer.PlayerId, 10)
	offerTeamId, _ := new(big.Int).SetString(inOffer.TeamId, 10)
	offerValidUntil, err := strconv.ParseInt(inOffer.ValidUntil, 10, 64)
	assert.NilError(t, err)

	hashOffer, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		offerValidUntil,
		offerPlayerId,
		extraPrice,
		dummyRnd,
		offerTeamId,
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

	assert.NilError(t, consumer.CreateOffer(tx, inOffer, *bc.Contracts))

	offerService := postgres.NewOfferService(tx)

	offer, err := offerService.OfferByRndPrice(inOffer.Rnd, inOffer.Price)
	assert.NilError(t, err)

	in := input.AcceptOfferInput{}
	in.ValidUntil = "999999999999"
	in.PlayerId = inOffer.PlayerId
	in.CurrencyId = inOffer.CurrencyId
	in.Price = inOffer.Price
	in.Rnd = inOffer.Rnd
	in.OfferId = graphql.ID(strconv.FormatInt(offer.ID, 10))

	playerId, _ := new(big.Int).SetString(in.PlayerId, 10)
	validUntil, err := strconv.ParseInt(in.ValidUntil, 10, 64)
	assert.NilError(t, err)
	hash, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		validUntil,
		playerId,
	)
	assert.Equal(t, hash.Hex(), "0xa22033e41e3e5e2f52acca342890f137d4dfa7007ab4741e9c289667a28067df")
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "c0ec61c2e32e37fe1bc1fbc8806c3f43461ac09cd4ac75b22c9625d9c80050da6b9d71f32800225c925f75a413620a8a29983c6cee138c5aab4593d7fbb897e51b")
	in.Signature = hex.EncodeToString(signature)

	assert.NilError(t, consumer.AcceptOffer(tx, in))
	auctionId, err := in.AuctionID()
	assert.NilError(t, err)

	service := postgres.NewAuctionService(tx)
	auction, err := service.Auction(string(auctionId))
	assert.NilError(t, err)
	assert.Assert(t, auction != nil)
	assert.Equal(t, auction.Seller, "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
	assert.Equal(t, auction.Price, int64(41234))

	offer, err = offerService.Offer(offer.ID)
	assert.NilError(t, err)
	assert.Equal(t, string(offer.State), "accepted")
	assert.Equal(t, offer.AuctionID, string(auctionId))

	bids, err := service.Bid().Bids(auction.ID)
	assert.Equal(t, string(bids[0].State), "accepted")
	assert.Equal(t, bids[0].Rnd, int64(inOffer.Rnd))
	assert.Equal(t, bids[0].ExtraPrice, int64(0))
	assert.Equal(t, string(bids[0].Signature), hex.EncodeToString(offerSignature))

}

func TestAcceptOfferWithExpiredOffer(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = "10000000"
	inOffer.PlayerId = "274877906940"
	inOffer.TeamId = "456678987944"
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = 4232
	inOffer.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f"
	offerPlayerId, _ := new(big.Int).SetString(inOffer.PlayerId, 10)
	offerTeamId, _ := new(big.Int).SetString(inOffer.TeamId, 10)
	offerValidUntil, err := strconv.ParseInt(inOffer.ValidUntil, 10, 64)
	assert.NilError(t, err)
	dummyRnd := int64(0)

	hashOffer, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		offerValidUntil,
		offerPlayerId,
		big.NewInt(0),
		big.NewInt(dummyRnd),
		offerTeamId,
		true,
	)
	assert.Equal(t, hashOffer.Hex(), "0x80ddf12ab28a6fb4a8ab17af2a81a7e251b5ca4d8aa1c97706218aa3782b7d1c")
	assert.NilError(t, err)
	offerPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	offerSignature, err := signer.Sign(hashOffer.Bytes(), offerPrivateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(offerSignature), "6ed548051674d96385ef4fc0e4dcd5e72697a125875eee0af85b94f0fe8c3dfd766dcde4fd4db9c566d444a8b698cf511969523229806c9c4a8e34c191c357681c")
	inOffer.Signature = hex.EncodeToString(offerSignature)

	assert.NilError(t, consumer.CreateOffer(tx, inOffer, *bc.Contracts))

	offerService := postgres.NewOfferService(tx)

	offer, err := offerService.OfferByRndPrice(inOffer.Rnd, inOffer.Price)
	assert.NilError(t, err)

	in := input.AcceptOfferInput{}
	in.ValidUntil = "999999999999"
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 4232
	in.OfferId = graphql.ID(strconv.FormatInt(offer.ID, 10))

	auctionPlayerId, _ := new(big.Int).SetString(in.PlayerId, 10)
	auctionValidUntil, err := strconv.ParseInt(in.ValidUntil, 10, 64)
	assert.NilError(t, err)
	hash, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		auctionValidUntil,
		auctionPlayerId,
	)
	assert.Equal(t, hash.Hex(), "0xf1d4501c5158a9018b1618ec4d471c66b663d8f6bffb6e70a0c6584f5c1ea94a")
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "381bf58829e11790830eab9924b123d1dbe96dd37b10112729d9d32d476c8d5762598042bb5d5fd63f668455aa3a2ce4e2632241865c26ababa231ad212b5f151b")
	in.Signature = hex.EncodeToString(signature)

	err = consumer.AcceptOffer(tx, in)
	assert.Error(t, err, "Associated Offer is expired")

}
