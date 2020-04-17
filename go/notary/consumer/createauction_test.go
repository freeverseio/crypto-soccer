package consumer_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/notary/consumer"
	"github.com/freeverseio/crypto-soccer/go/notary/producer/gql/input"
	"github.com/freeverseio/crypto-soccer/go/notary/signer"
	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"gotest.tools/assert"
)

func TestCreateAuction(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	in := input.CreateAuctionInput{}
	in.ValidUntil = fmt.Sprintf("%v", time.Now().Unix()+1000)
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
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	in.Signature = hex.EncodeToString(signature)

	assert.NilError(t, consumer.CreateAuction(tx, in))
	id, err := in.ID()
	assert.NilError(t, err)
	auction, err := storage.AuctionByID(tx, string(id))
	assert.NilError(t, err)
	assert.Assert(t, auction != nil)
	assert.Equal(t, auction.Seller, "0x83A909262608c650BD9b0ae06E29D90D0F67aC5e")
}
