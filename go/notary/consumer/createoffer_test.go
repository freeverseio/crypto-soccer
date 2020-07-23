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

func TestCreateOffer(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	in := input.CreateOfferInput{}
	in.ValidUntil = "999999999999"
	in.PlayerId = "274877906940"
	in.TeamId = "456678987944"
	in.CurrencyId = 1
	in.Price = 41234
	in.Rnd = 4232
	in.Seller = "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f"
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
	assert.Equal(t, hash.Hex(), "0xe194d576ea5dff0e13e4f9d9d2aa4f5fb06af68732fd0e2106c82a8e7949ef19")
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "83df48b4f3be03a020690b7af318f9cf4005a874da05631443dc833c3644c38865383b54d2556e8c50de6a15d65d88c4214ea849c0867b91cdb946de4b24bee11c")
	in.Signature = hex.EncodeToString(signature)

	assert.NilError(t, consumer.CreateOffer(tx, in))
	id, err := in.ID()
	assert.NilError(t, err)

	service := postgres.NewOfferService(tx)
	offer, err := service.Offer(string(id))
	assert.NilError(t, err)
	assert.Assert(t, offer != nil)
	assert.Equal(t, offer.Seller, "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f")
}
