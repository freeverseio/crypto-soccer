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
		big.NewInt(0),
		big.NewInt(int64(inOffer.Rnd)),
		offerTeamId,
		true,
	)
	assert.Equal(t, hashOffer.Hex(), "0x5c3817ae7930907579b9694a5f5439906c1695a6985e772f982ff7fea2f9ae7e")
	assert.NilError(t, err)
	offerPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	offerSignature, err := signer.Sign(hashOffer.Bytes(), offerPrivateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(offerSignature), "030e2a64488d1fe9150cf0ca9c85ca275f3f867c905a4e325eb1715d9d39621207cd3001c929f00ed1b56370d21def0b1179a7b7c1d27602f80206636fc1b2961c")
	inOffer.Signature = hex.EncodeToString(offerSignature)

	assert.NilError(t, consumer.CreateOffer(tx, inOffer, *bc.Contracts))
	idOffer, err := inOffer.ID(*bc.Contracts)
	assert.NilError(t, err)

	offerService := postgres.NewOfferService(tx)

	offer, err := offerService.Offer(string(idOffer))
	assert.NilError(t, err)
	assert.Equal(t, offer.ID, string(idOffer))

	in := input.AcceptOfferInput{}
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

	assert.NilError(t, consumer.AcceptOffer(tx, in))
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

	hashOffer, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		offerValidUntil,
		offerPlayerId,
		big.NewInt(0),
		big.NewInt(int64(inOffer.Rnd)),
		offerTeamId,
		true,
	)
	assert.Equal(t, hashOffer.Hex(), "0x77ea6c76e92738e3dafbb9bc9cfd9e726671a474eb7b659303bf80dc1b14ebe0")
	assert.NilError(t, err)
	offerPrivateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	offerSignature, err := signer.Sign(hashOffer.Bytes(), offerPrivateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(offerSignature), "7f0c39cf5496deab81f6e597c2689853817cfce0a490bd3355bb636beaaffdcc6de5d20949d1cc6561fc069c3ba0708924d6aae11952f9fc961bd98858f931251c")
	inOffer.Signature = hex.EncodeToString(offerSignature)

	assert.NilError(t, consumer.CreateOffer(tx, inOffer, *bc.Contracts))
	idOffer, err := inOffer.ID(*bc.Contracts)
	assert.NilError(t, err)

	offerService := postgres.NewOfferService(tx)

	offer, err := offerService.Offer(string(idOffer))
	assert.NilError(t, err)
	assert.Equal(t, offer.ID, string(idOffer))

	in := input.AcceptOfferInput{}
	in.ValidUntil = "999999999999"
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 4232
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
