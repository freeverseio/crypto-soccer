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

	batch, err := service.CreateBatch(teamId, now)
	assert.NilError(t, err)
	assert.Equal(t, len(batch), 25)

	playerId = string(batch[1].PlayerId())
	wp, err = service.GetWorldPlayer(playerId, teamId, now)
	assert.NilError(t, err)
	golden.Assert(t, dump.Sdump(wp), t.Name()+".golden")
}
