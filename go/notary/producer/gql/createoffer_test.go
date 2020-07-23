package gql_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

//TODOO
func TestCreateOfferReturnTheSignature(t *testing.T) {
	// step 1: offerer creates a random number
	offererRnd := int32(42321)

	// step 2: the seller puts for sale (freezes) his player using the buyer RND number
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, db)

	in := input.CreateAuctionInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	in.PlayerId = "2748779069494"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = offererRnd

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

	// Step 3: now that player is frozen, if noone else bids, we can send directly the offer sign msg to the BC

	inOffer := input.CreateOfferInput{}
	inOffer.ValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	inOffer.PlayerId = in.PlayerId
	inOffer.CurrencyId = in.CurrencyId
	inOffer.Price = in.Price
	inOffer.Rnd = offererRnd

	inOffer.TeamId = "2748779069441"
	teamId, _ := new(big.Int).SetString(in.TeamId, 10)

	hash, err := signer.HashOfferMessage(
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		validUntil,
		playerId,
		teamId,
	)
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	// buyer, err := crypto.HexToECDSA("FE158D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	// signature, err := signer.Sign(hash.Bytes(), buyer)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)
	_, err = r.CreateAuction(struct{ Input input.CreateAuctionInput }{in})
	assert.NilError(t, err)

	// id, err := r.CreateOffer(struct{ Input input.CreateOfferInput }{in})
	// assert.NilError(t, err)
	// id2, err := in.ID()
	// assert.NilError(t, err)
	// assert.Equal(t, id, id2)
}
