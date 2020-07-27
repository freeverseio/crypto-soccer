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
	hash, err := signer.HashBidMessage(
		bc.Contracts.Market,
		uint8(in.CurrencyId),
		big.NewInt(int64(in.Price)),
		big.NewInt(int64(in.Rnd)),
		validUntil,
		playerId,
		big.NewInt(0),
		big.NewInt(int64(in.Rnd)),
		teamId,
		true,
	)
	assert.Equal(t, hash.Hex(), "0x5c3817ae7930907579b9694a5f5439906c1695a6985e772f982ff7fea2f9ae7e")
	assert.NilError(t, err)
	privateKey, err := crypto.HexToECDSA("FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85")
	assert.NilError(t, err)
	signature, err := signer.Sign(hash.Bytes(), privateKey)
	assert.NilError(t, err)
	assert.Equal(t, hex.EncodeToString(signature), "030e2a64488d1fe9150cf0ca9c85ca275f3f867c905a4e325eb1715d9d39621207cd3001c929f00ed1b56370d21def0b1179a7b7c1d27602f80206636fc1b2961c")
	in.Signature = hex.EncodeToString(signature)

	assert.NilError(t, consumer.CreateOffer(tx, in, *bc.Contracts))
	id, err := in.ID(*bc.Contracts)
	assert.NilError(t, err)

	service := postgres.NewOfferService(tx)
	offer, err := service.Offer(string(id))
	assert.NilError(t, err)
	assert.Assert(t, offer != nil)
	assert.Equal(t, offer.Seller, "0x83A909262608c650BD9b0ae06E29D90D0F67aC5f")
}
