package gql_test

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"gotest.tools/assert"
)

func TestCreateAuctionReturnTheSignature(t *testing.T) {
	ch := make(chan interface{}, 10)
	r := gql.NewResolver(ch, *bc.Contracts)

	in := input.CreateAuctionInput{}
	in.ValidUntil = "5453636457457456"
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
	signature, err := signer.Sign(hash.Bytes(), bc.Owner)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "e04252b5dc5efdefb97d53e135add1a7cd3e15279dbe6519847d86d4fd52fbb5777e4880c9c8166f27490d73c7e687d41089b26e2572f19c17f1956ddef6b49c1c")
	in.Signature = hex.EncodeToString(signature)

	id, err := r.CreateAuction(struct{ Input input.CreateAuctionInput }{in})
	assert.NilError(t, err)
	id2, err := in.ID()
	assert.NilError(t, err)
	assert.Equal(t, id, id2)
}
