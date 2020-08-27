package consumer_test

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/graph-gophers/graphql-go"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"
	"gotest.tools/assert"
)

func TestAcceptOffer(t *testing.T) {
	service := postgres.NewStorageService(db)
	assert.NilError(t, service.Begin())
	defer service.Rollback()
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
	inOffer.BuyerTeamId = ofererTeamId.String()
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

	assert.NilError(t, consumer.CreateOffer(service, inOffer, *bc.Contracts))

	offer, err := service.OfferByRndPrice(inOffer.Rnd, inOffer.Price)
	assert.NilError(t, err)

	in := input.AcceptOfferInput{}
	in.ValidUntil = auctionValidUntil.String()
	in.PlayerId = inOffer.PlayerId
	in.CurrencyId = inOffer.CurrencyId
	in.Price = inOffer.Price
	in.Rnd = inOffer.Rnd
	in.OfferId = graphql.ID(offer.ID)

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

	assert.NilError(t, consumer.AcceptOffer(service, in))
	auctionId, err := in.AuctionID()
	assert.NilError(t, err)

	auction, err := service.Auction(string(auctionId))
	assert.NilError(t, err)
	assert.Assert(t, auction != nil)
	assert.Equal(t, auction.Seller, "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
	assert.Equal(t, auction.Price, int64(41234))

	offer, err = service.Offer(offer.ID)
	assert.NilError(t, err)
	assert.Equal(t, string(offer.State), "accepted")
	assert.Equal(t, offer.AuctionID, string(auctionId))

	bids, err := service.Bids(auction.ID)
	assert.Equal(t, string(bids[0].State), "accepted")
	assert.Equal(t, bids[0].Rnd, int64(inOffer.Rnd))
	assert.Equal(t, bids[0].ExtraPrice, int64(0))
	assert.Equal(t, string(bids[0].Signature), hex.EncodeToString(offerSignature))

}

func TestAcceptOfferWithExpiredOffer(t *testing.T) {
	service := postgres.NewStorageService(db)
	assert.NilError(t, service.Begin())
	defer service.Rollback()

	// in this example the offer expired 5 min ago
	offerValidUntil := big.NewInt(time.Now().Unix() - 5*60)
	auctionValidUntil := big.NewInt(time.Now().Unix() + 48*3600)
	playerId := big.NewInt(274877906940)
	ofererTeamId := big.NewInt(456678987944)

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = offerValidUntil.String()
	inOffer.PlayerId = playerId.String()
	inOffer.BuyerTeamId = ofererTeamId.String()
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
	// We cannot compare hashes because validUntil is different in every test (it uses time.now())
	assert.NilError(t, err)
	offerPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	offerSignature, err := signer.Sign(hashOffer.Bytes(), offerPrivateKey)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(offerSignature)

	assert.NilError(t, consumer.CreateOffer(service, inOffer, *bc.Contracts))

	offer, err := service.OfferByRndPrice(inOffer.Rnd, inOffer.Price)
	assert.NilError(t, err)

	in := input.AcceptOfferInput{}
	in.ValidUntil = auctionValidUntil.String()
	in.PlayerId = inOffer.PlayerId
	in.CurrencyId = inOffer.CurrencyId
	in.Price = inOffer.Price
	in.Rnd = inOffer.Rnd
	in.OfferId = graphql.ID(offer.ID)

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

	err = consumer.AcceptOffer(service, in)
	assert.Error(t, err, "Associated Offer is expired")

}

func TestAcceptOfferWithNonExpiredOffer(t *testing.T) {
	service := postgres.NewStorageService(db)
	assert.NilError(t, service.Begin())
	defer service.Rollback()

	// in this example the offer expires in 5 min
	offerValidUntil := big.NewInt(time.Now().Unix() + 5*60)
	auctionValidUntil := big.NewInt(time.Now().Unix() + 48*3600)
	playerId := big.NewInt(274877906940)
	ofererTeamId := big.NewInt(456678987944)

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = offerValidUntil.String()
	inOffer.PlayerId = playerId.String()
	inOffer.BuyerTeamId = ofererTeamId.String()
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
	// We cannot compare hashes because validUntil is different in every test (it uses time.now())
	assert.NilError(t, err)
	offerPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	offerSignature, err := signer.Sign(hashOffer.Bytes(), offerPrivateKey)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(offerSignature)

	assert.NilError(t, consumer.CreateOffer(service, inOffer, *bc.Contracts))

	offer, err := service.OfferByRndPrice(inOffer.Rnd, inOffer.Price)
	assert.NilError(t, err)

	in := input.AcceptOfferInput{}
	in.ValidUntil = auctionValidUntil.String()
	in.PlayerId = inOffer.PlayerId
	in.CurrencyId = inOffer.CurrencyId
	in.Price = inOffer.Price
	in.Rnd = inOffer.Rnd
	in.OfferId = graphql.ID(offer.ID)

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

	err = consumer.AcceptOffer(service, in)
	assert.NilError(t, err)
}

func TestAcceptOfferWithNonHighestOffer(t *testing.T) {
	service := postgres.NewStorageService(db)
	assert.NilError(t, service.Begin())
	defer service.Rollback()

	playerId := big.NewInt(274877906940)
	extraPrice := big.NewInt(0)
	dummyRnd := big.NewInt(0)

	// Creating lower offer
	lowerOfferValidUntil := big.NewInt(time.Now().Unix() + 5*60)
	lowerAuctionValidUntil := big.NewInt(time.Now().Unix() + 48*3600)
	lowerOfererTeamId := big.NewInt(456678987944)

	inLowerOffer := input.CreateOfferInput{}
	inLowerOffer.ValidUntil = lowerOfferValidUntil.String()
	inLowerOffer.PlayerId = playerId.String()
	inLowerOffer.BuyerTeamId = lowerOfererTeamId.String()
	inLowerOffer.CurrencyId = 1
	inLowerOffer.Price = 1000
	inLowerOffer.Rnd = 4232
	inLowerOffer.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f"

	hashLowerOffer, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inLowerOffer.CurrencyId),
		big.NewInt(int64(inLowerOffer.Price)),
		big.NewInt(int64(inLowerOffer.Rnd)),
		lowerOfferValidUntil.Int64(),
		playerId,
		extraPrice,
		dummyRnd,
		lowerOfererTeamId,
		true,
	)
	// We cannot compare hashes because validUntil is different in every test (it uses time.now())
	assert.NilError(t, err)
	lowerOfferPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	lowerOfferSignature, err := signer.Sign(hashLowerOffer.Bytes(), lowerOfferPrivateKey)
	assert.NilError(t, err)
	inLowerOffer.Signature = hex.EncodeToString(lowerOfferSignature)

	assert.NilError(t, consumer.CreateOffer(service, inLowerOffer, *bc.Contracts))

	lowerOffer, err := service.OfferByRndPrice(inLowerOffer.Rnd, inLowerOffer.Price)
	assert.NilError(t, err)

	// Creating highest offer
	highestOfferValidUntil := big.NewInt(time.Now().Unix() + 5*60)
	highestAuctionValidUntil := big.NewInt(time.Now().Unix() + 48*3600)
	highesetOfererTeamId := big.NewInt(456678987944)

	inHighestOffer := input.CreateOfferInput{}
	inHighestOffer.ValidUntil = highestOfferValidUntil.String()
	inHighestOffer.PlayerId = playerId.String()
	inHighestOffer.BuyerTeamId = highesetOfererTeamId.String()
	inHighestOffer.CurrencyId = 1
	inHighestOffer.Price = 2000
	inHighestOffer.Rnd = 4232
	inHighestOffer.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f"
	assert.NilError(t, err)

	hashHighestOffer, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inHighestOffer.CurrencyId),
		big.NewInt(int64(inHighestOffer.Price)),
		big.NewInt(int64(inHighestOffer.Rnd)),
		highestOfferValidUntil.Int64(),
		playerId,
		extraPrice,
		dummyRnd,
		highesetOfererTeamId,
		true,
	)
	// We cannot compare hashes because validUntil is different in every test (it uses time.now())
	assert.NilError(t, err)
	highestOfferPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	highestOfferSignature, err := signer.Sign(hashHighestOffer.Bytes(), highestOfferPrivateKey)
	assert.NilError(t, err)
	inHighestOffer.Signature = hex.EncodeToString(highestOfferSignature)

	assert.NilError(t, consumer.CreateOffer(service, inHighestOffer, *bc.Contracts))

	higherOffer, err := service.OfferByRndPrice(inHighestOffer.Rnd, inHighestOffer.Price)
	assert.NilError(t, err)
	_, err = service.OfferByRndPrice(inHighestOffer.Rnd, inHighestOffer.Price)
	assert.NilError(t, err)

	// Accepting lower offer
	in := input.AcceptOfferInput{}
	in.ValidUntil = lowerAuctionValidUntil.String()
	in.PlayerId = inLowerOffer.PlayerId
	in.CurrencyId = inLowerOffer.CurrencyId
	in.Price = inLowerOffer.Price
	in.Rnd = inLowerOffer.Rnd
	in.OfferId = graphql.ID(lowerOffer.ID)

	assert.NilError(t, err)
	hash, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		lowerAuctionValidUntil.Int64(),
		playerId,
	)
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	err = consumer.AcceptOffer(service, in)
	assert.Error(t, err, "You can only accept highest offer")

	// Accepting hihgest offer
	inHigher := input.AcceptOfferInput{}
	inHigher.ValidUntil = highestAuctionValidUntil.String()
	inHigher.PlayerId = inHighestOffer.PlayerId
	inHigher.CurrencyId = inHighestOffer.CurrencyId
	inHigher.Price = inHighestOffer.Price
	inHigher.Rnd = inHighestOffer.Rnd
	inHigher.OfferId = graphql.ID(higherOffer.ID)

	hashHigher, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		highestAuctionValidUntil.Int64(),
		playerId,
	)
	signature, err = signer.Sign(hashHigher.Bytes(), privateKey)
	assert.NilError(t, err)
	inHigher.Signature = hex.EncodeToString(signature)

	err = consumer.AcceptOffer(service, inHigher)
	assert.NilError(t, err)

	// Accepting unexistent offer
	inUnexistent := input.AcceptOfferInput{}
	inUnexistent.ValidUntil = big.NewInt(999999999999).String()
	inUnexistent.PlayerId = playerId.String()
	inUnexistent.CurrencyId = 1
	inUnexistent.Price = 100
	inUnexistent.Rnd = 45654
	inUnexistent.OfferId = graphql.ID("23442566")

	hashUnexistent, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		highestAuctionValidUntil.Int64(),
		playerId,
	)
	signature, err = signer.Sign(hashUnexistent.Bytes(), privateKey)
	assert.NilError(t, err)
	inUnexistent.Signature = hex.EncodeToString(signature)

	err = consumer.AcceptOffer(service, inUnexistent)
	assert.Error(t, err, "You can only accept highest offer")
}

func TestAcceptUnexistentOffers(t *testing.T) {
	service := postgres.NewStorageService(db)
	assert.NilError(t, service.Begin())
	defer service.Rollback()

	auctionValidUntil := big.NewInt(999999999999 + 3600*24*2)
	playerId := big.NewInt(274877906940)

	in := input.AcceptOfferInput{}
	in.ValidUntil = big.NewInt(999999999999).String()
	in.PlayerId = playerId.String()
	in.CurrencyId = 1
	in.Price = 100
	in.Rnd = 45654
	in.OfferId = graphql.ID("23442566")

	hash, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		auctionValidUntil.Int64(),
		playerId,
	)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	err = consumer.AcceptOffer(service, in)
	assert.Error(t, err, "There are no offers for this playerId")

}
