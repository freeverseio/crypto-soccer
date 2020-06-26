package worldplayer_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/notary/worldplayer"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestWorldPlayerService(t *testing.T) {
	now := int64(1554940800) // first second of a week
	teamId := "274877906944"

	service := worldplayer.NewWorldPlayerService(*bc.Contracts, namesdb)

	playerId := "0"
	wp, err := service.GetWorldPlayer(playerId, teamId, now)
	assert.NilError(t, err)
	assert.Assert(t, wp == nil)

	// given this teamId and time, you get 32 players because noone shows up in tier3
	batch, err := service.CreateBatch(teamId, now)
	assert.NilError(t, err)
	assert.Equal(t, len(batch), 32)

	playerId = string(batch[1].PlayerId())
	wp, err = service.GetWorldPlayer(playerId, teamId, now)
	assert.NilError(t, err)
	golden.Assert(t, dump.Sdump(wp), t.Name()+".golden")

	// changing teamId or time can lead to existing tier3 players
	// First, change teamId:
	teamId2 := "274877906946"
	batch, err = service.CreateBatch(teamId2, now)
	assert.NilError(t, err)
	assert.Equal(t, len(batch), 33)

	// Second, change "now":
	now2 := now + now/4
	batch, err = service.CreateBatch(teamId, now2)
	assert.NilError(t, err)
	assert.Equal(t, len(batch), 33)
}
