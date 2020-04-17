package gql_test

import (
	"encoding/hex"
	"fmt"
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

func TestCreateAuctionReturnTheSignature(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts, db)

	in := input.CreateAuctionInput{}
	in.ValidUntil = fmt.Sprintf("%v", time.Now().Unix()+1000)
	in.PlayerId = "274877906944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 42321

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
	assert.NilError(t, err)
	signature, err := crypto.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	id, err := r.CreateAuction(struct{ Input input.CreateAuctionInput }{in})
	assert.NilError(t, err)
	id2, err := in.ID()
	assert.NilError(t, err)
	assert.Equal(t, id, id2)
}
