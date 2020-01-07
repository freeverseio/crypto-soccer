package match_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
	"gotest.tools/assert"
)

func TestCreateDummyPlayer(t *testing.T) {
	t.Parallel()
	defence := uint16(0)
	speed := uint16(0)
	endurance := uint16(0)
	pass := uint16(0)
	shoot := uint16(0)
	player := match.CreateDummyPlayer(t, bc.Contracts, defence, speed, endurance, pass, shoot)
	value, err := player.Defence(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, value, defence)
	assert.Equal(t, player.Skills().String(), "6368953449211795048194580334409608269205078016")
	birth, err := player.Birth(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, birth.UTC().String(), "2016-10-04 00:00:00 +0000 UTC")
}
