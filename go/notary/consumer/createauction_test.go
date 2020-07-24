package consumer_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/storage/postgres"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

func TestCreateAuction(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	in := input.CreateAuctionInput{}
	in.ValidUntil = "999999999999"
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 4232
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
	assert.Equal(t, hash.Hex(), "0xf1d4501c5158a9018b1618ec4d471c66b663d8f6bffb6e70a0c6584f5c1ea94a")
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "381bf58829e11790830eab9924b123d1dbe96dd37b10112729d9d32d476c8d5762598042bb5d5fd63f668455aa3a2ce4e2632241865c26ababa231ad212b5f151b")
	in.Signature = hex.EncodeToString(signature)

	assert.NilError(t, consumer.CreateAuction(tx, in))
	id, err := in.ID()
	assert.NilError(t, err)

	service := postgres.NewAuctionService(tx)
	auction, err := service.Auction(string(id))
	assert.NilError(t, err)
	assert.Assert(t, auction != nil)
	assert.Equal(t, auction.Seller, "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
}

func TestCreateAuctionWithOffer(t *testing.T) {
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
	offerPlayerId, _ := new(big.Int).SetString(inOffer.PlayerId, 10)
	offerTeamId, _ := new(big.Int).SetString(inOffer.TeamId, 10)
	offerValidUntil, err := strconv.ParseInt(inOffer.ValidUntil, 10, 64)
	assert.NilError(t, err)
	hashOffer, err := signer.HashOfferMessage(
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		offerValidUntil,
		offerPlayerId,
		offerTeamId,
	)
	assert.Equal(t, hashOffer.Hex(), "0xe194d576ea5dff0e13e4f9d9d2aa4f5fb06af68732fd0e2106c82a8e7949ef19")
	assert.NilError(t, err)
	offerPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	offerSignature, err := signer.Sign(hashOffer.Bytes(), offerPrivateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(offerSignature), "83df48b4f3be03a020690b7af318f9cf4005a874da05631443dc833c3644c38865383b54d2556e8c50de6a15d65d88c4214ea849c0867b91cdb946de4b24bee11c")
	inOffer.Signature = hex.EncodeToString(offerSignature)

	assert.NilError(t, consumer.CreateOffer(tx, inOffer))
	idOffer, err := inOffer.ID()
	assert.NilError(t, err)

	offerService := postgres.NewOfferService(tx)

	offer, err := offerService.Offer(string(idOffer))
	assert.NilError(t, err)
	assert.Equal(t, offer.ID, string(idOffer))

	in := input.CreateAuctionInput{}
	in.ValidUntil = "999999999999"
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 4232
	in.OfferId = string(idOffer)
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
	assert.Equal(t, hash.Hex(), "0xf1d4501c5158a9018b1618ec4d471c66b663d8f6bffb6e70a0c6584f5c1ea94a")
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "381bf58829e11790830eab9924b123d1dbe96dd37b10112729d9d32d476c8d5762598042bb5d5fd63f668455aa3a2ce4e2632241865c26ababa231ad212b5f151b")
	in.Signature = hex.EncodeToString(signature)

	assert.NilError(t, consumer.CreateAuction(tx, in))
	id, err := in.ID()
	assert.NilError(t, err)

	service := postgres.NewAuctionService(tx)
	auction, err := service.Auction(string(id))
	assert.NilError(t, err)
	assert.Assert(t, auction != nil)
	assert.Equal(t, auction.Seller, "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")

	offer, err = offerService.Offer(string(idOffer))
	assert.NilError(t, err)
	assert.Equal(t, string(offer.State), "accepted")
	assert.Equal(t, offer.AuctionID, string(id))
}
