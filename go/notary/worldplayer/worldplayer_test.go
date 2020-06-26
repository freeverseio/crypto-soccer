package worldplayer_test

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/freeverseio/crypto-soccer/go/helper"
	"github.com/freeverseio/crypto-soccer/go/notary/worldplayer"
	"github.com/freeverseio/crypto-soccer/go/testutils"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestGetWorldPlayersDeterministicResult(t *testing.T) {
	now := int64(3600 * 24 * 7 * 2571) // first second of a week, and of a day
	teamId := "274877906944"

	service := worldplayer.NewWorldPlayerService(*bc.Contracts, namesdb)
	players0, err := service.CreateBatch(teamId, now)
	assert.NilError(t, err)
	players1, err := service.CreateBatch(teamId, now)
	assert.NilError(t, err)

	assert.Equal(t, len(players0), 32)
	assert.Equal(t, len(players0), len(players1))
	for i := range players0 {
		assert.Equal(t, *players0[i], *players1[i])
	}

	golden.Assert(t, dump.Sdump(players0), t.Name()+".golden")
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
	assert.Equal(t, len(players), 32)

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
	assert.Equal(t, len(players), 31)
}

func TestWorldPlayerServiceTimeDependentDiscontinuity(t *testing.T) {
	now := int64(3600 * 24 * 7 * 2571) // first second of a week, and of a day
	teamId := "274877906944"
	halfDay := int64(3600 * 12)

	// first, check that if created at now-1 or now-2, the offering time would be 12 hours ago,
	// and hence, the validUntil would be exactly now
	service := worldplayer.NewWorldPlayerService(*bc.Contracts, namesdb)
	players, err := service.CreateBatch(teamId, now-2)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 32)
	assert.Equal(t, string(players[0].PlayerId()), "25726027164631419126757699466463075271849876153057056752928093")
	expectedValidUntil := int64(1554940800)
	assert.Equal(t, players[0].ValidUntil(), strconv.FormatUint(uint64(expectedValidUntil), 10))
	assert.Equal(t, expectedValidUntil-now == 0, true)

	players, err = service.CreateBatch(teamId, now-1)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 32)
	assert.Equal(t, string(players[0].PlayerId()), "25726027164631419126757699466463075271849876153057056752928093")
	expectedValidUntil = int64(1554940800)
	assert.Equal(t, players[0].ValidUntil(), strconv.FormatUint(uint64(expectedValidUntil), 10))
	assert.Equal(t, expectedValidUntil-now == 0, true)

	// second, check that if created at now or now+1, the offering time would be exactly now,
	// and hence, the validUntil would be exactly 12 hours in the future
	players, err = service.CreateBatch(teamId, now)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 32)
	assert.Equal(t, string(players[0].PlayerId()), "25726027164631419126757699466464342922450104382458553456133469")
	expectedValidUntil = int64(1554984000)
	assert.Equal(t, players[0].ValidUntil(), strconv.FormatUint(uint64(expectedValidUntil), 10))
	assert.Equal(t, expectedValidUntil-now == halfDay, true)

	players, err = service.CreateBatch(teamId, now+1)
	assert.NilError(t, err)
	assert.Equal(t, len(players), 32)
	assert.Equal(t, string(players[0].PlayerId()), "25726027164631419126757699466464342922450104382458553456133469")
	expectedValidUntil = int64(1554984000)
	assert.Equal(t, players[0].ValidUntil(), strconv.FormatUint(uint64(expectedValidUntil), 10))
	assert.Equal(t, expectedValidUntil-now == halfDay, true)
}
