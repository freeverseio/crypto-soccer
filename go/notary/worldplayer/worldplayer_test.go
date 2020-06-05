package worldplayer_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/worldplayer"
	"github.com/freeverseio/crypto-soccer/go/testutils"
	"gotest.tools/assert"
)

func TestGetWorldPlayersDeterministicResult(t *testing.T) {
	now := int64(1554940800) // first second of a week
	teamId := "274877906944"

	service := worldplayer.NewWorldPlayerService(*bc.Contracts, namesdb)
	players0, err := service.CreateBatch(teamId, now)
	assert.NilError(t, err)
	players1, err := service.CreateBatch(teamId, now)
	assert.NilError(t, err)

	assert.Equal(t, len(players0), 30)
	assert.Equal(t, len(players0), len(players1))
	for i := range players0 {
		assert.Equal(t, *players0[i], *players1[i])
	}
	assert.Equal(t, players0[0].Race() == "", false)
	assert.Equal(t, players0[0].CountryOfBirth() == "", false)
}

func TestGetWorldPlayersOfSoldPlayer(t *testing.T) {
	bc, err := testutils.NewBlockchain()
	assert.NilError(t, err)
	timezoneIdx := uint8(1)
	countryIdx := big.NewInt(0)
	bc.Contracts.Assets.TransferFirstBotToAddr(
		bind.NewKeyedTransactor(bc.Owner),
		timezoneIdx,
		countryIdx,
		crypto.PubkeyToAddress(bc.Owner.PublicKey),
	)

	now := int64(1554940801) // first second of a week
	teamId := "274877906944"

	service := worldplayer.NewWorldPlayerService(*bc.Contracts, namesdb)
	players, err := service.CreateBatch(teamId, now)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, players[0].Race() == "", false)
	assert.Equal(t, players[0].CountryOfBirth() == "", false)

	player0Id, _ := new(big.Int).SetString(string(players[0].PlayerId()), 10)
	targetTeamId, _ := new(big.Int).SetString(teamId, 10)
	tx, err := bc.Contracts.Market.TransferBuyNowPlayer(
		bind.NewKeyedTransactor(bc.Owner),
		player0Id,
		targetTeamId,
	)
	assert.NilError(t, err)
	_, err = helper.WaitReceipt(bc.Client, tx, 60)
	assert.NilError(t, err)

	players, err = service.CreateBatch(teamId, now)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 29)
	assert.Equal(t, players[0].Race() == "", false)
	assert.Equal(t, players[0].CountryOfBirth() == "", false)
}

func TestWorldPlayerServiceTimeDependentDiscontinuity(t *testing.T) {
	now := int64(1554940800) // first second of a week
	teamId := "274877906944"

	service := worldplayer.NewWorldPlayerService(*bc.Contracts, namesdb)
	players, err := service.CreateBatch(teamId, now-2)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "25725015717574447640111183123578655469576475906857003090183425")
	assert.Equal(t, players[0].ValidUntil(), "1554940800")

	players, err = service.CreateBatch(teamId, now-1)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "25725015717574447640111183123578655469576475906857003090183425")
	assert.Equal(t, players[0].ValidUntil(), "1554940800")

	players, err = service.CreateBatch(teamId, now)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "25725015717574447640111183123579923120176704136258499793388801")
	assert.Equal(t, players[0].ValidUntil(), "1555545600")

	players, err = service.CreateBatch(teamId, now+1)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 30)
	assert.Equal(t, string(players[0].PlayerId()), "25725015717574447640111183123579923120176704136258499793388801")
	assert.Equal(t, players[0].ValidUntil(), "1555545600")
}
