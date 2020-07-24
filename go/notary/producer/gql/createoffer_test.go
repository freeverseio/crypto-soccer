package gql_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

func TestCreateOffer(t *testing.T) {
	offererRnd := int32(42321)
	offerValidUntil := time.Now().Unix() + 100

	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, db)

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = strconv.FormatInt(offerValidUntil, 10)
	inOffer.PlayerId = "274877906944"
	inOffer.CurrencyId = 1
	inOffer.Price = 41234
	inOffer.Rnd = offererRnd
	inOffer.TeamId = "274877906945"
	teamId, _ := new(big.Int).SetString(inOffer.TeamId, 10)
	playerId, _ := new(big.Int).SetString(inOffer.PlayerId, 10)
	validUntil, err := strconv.ParseInt(inOffer.ValidUntil, 10, 64)

	hashOffer, err := signer.HashOfferMessage(
		uint8(inOffer.CurrencyId),
		big.NewInt(int64(inOffer.Price)),
		big.NewInt(int64(inOffer.Rnd)),
		validUntil,
		playerId,
		teamId,
	)
	assert.NilError(t, err)
	buyerKey, err := crypto.HexToECDSA("3B878F7892FBBFA30C8AED1DF317C19B853685E707C2CF0EE1927DC516060A54")
	signatureOffer, err := signer.Sign(hashOffer.Bytes(), buyerKey)
	assert.NilError(t, err)
	inOffer.Signature = hex.EncodeToString(signatureOffer)
	id, err := r.CreateOffer(struct{ Input input.CreateOfferInput }{inOffer})
	assert.NilError(t, err)
	id2, err := inOffer.ID()
	assert.NilError(t, err)
	assert.Equal(t, id, id2)

	in := input.CreateAuctionInput{}
	in.ValidUntil = strconv.FormatInt(offerValidUntil+1000, 10)
	in.PlayerId = inOffer.PlayerId
	in.CurrencyId = inOffer.CurrencyId
	in.Price = inOffer.Price
	in.Rnd = offererRnd
	auctionId, err := in.ID()

	validUntilAuction, err := strconv.ParseInt(in.ValidUntil, 10, 64)
	assert.NilError(t, err)
	hash, err := signer.HashSellMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		validUntilAuction,
		playerId,
	)

	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	_, err = r.CreateAuction(struct{ Input input.CreateAuctionInput }{in})
	assert.NilError(t, err)

	inBid := input.CreateBidInput{}
	inBid.AuctionId = auctionId
	inBid.ExtraPrice = 0
	inBid.Rnd = offererRnd
	inBid.TeamId = inOffer.TeamId
	inBid.Signature = inOffer.Signature

	_, err = r.CreateBid(struct{ Input input.CreateBidInput }{inBid})
	assert.NilError(t, err)
}
