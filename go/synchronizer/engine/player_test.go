package engine_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/engine"
	"gotest.tools/assert"
)

func TestNullPlayer(t *testing.T) {
	t.Parallel()
	p := engine.NewPlayer()
	assert.Equal(t, p.DumpState(), "skills: 0")
}

func TestPlayerToStorage(t *testing.T) {
	player := engine.NewPlayer()
	player.SetSkills(SkillsFromString("14606253788909032162646379450304996475079674564248175"))
	sto, err := player.ToStorage(*bc.Contracts)
	assert.NilError(t, err)
	assert.Equal(t, sto.Defence, uint64(955))
	assert.Equal(t, sto.Speed, uint64(889))
	assert.Equal(t, sto.Pass, uint64(1076))
	assert.Equal(t, sto.Endurance, uint64(1454))
	assert.Equal(t, sto.Shoot, uint64(623))
}
