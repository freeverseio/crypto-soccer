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

func TestCreateOfferReturnTheSignature(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, namesdb, googleCredentials, db)

	in := input.CreateOfferInput{}
	in.ValidUntil = strconv.FormatInt(time.Now().Unix()+100, 10)
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321
	in.TeamId = "34234234"

	playerId, _ := new(big.Int).SetString(in.PlayerId, 10)
	teamId, _ := new(big.Int).SetString(in.TeamId, 10)
	validUntil, err := strconv.ParseInt(in.ValidUntil, 10, 64)
	assert.NilError(t, err)
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
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	id, err := r.CreateOffer(struct{ Input input.CreateOfferInput }{in})
	assert.NilError(t, err)
	id2, err := in.ID()
	assert.NilError(t, err)
	assert.Equal(t, id, id2)
}
