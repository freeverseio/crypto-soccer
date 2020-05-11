package engine_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/engine"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestNullPlayer(t *testing.T) {
	t.Parallel()
	p := engine.NewPlayer()
	golden.Assert(t, dump.Sdump(p), t.Name()+".golden")
}

func TestPlayerToStorage(t *testing.T) {
	player := engine.NewPlayer()
	player.SetSkills(*bc.Contracts, SkillsFromString(t, "14606253788909032162646379450304996475079674564248175"))
	assert.Equal(t, player.Defence, uint64(23264))
	assert.Equal(t, player.Speed, uint64(213047))
	assert.Equal(t, player.Pass, uint64(244484))
	assert.Equal(t, player.Endurance, uint64(16970))
	assert.Equal(t, player.Shoot, uint64(590447))
}

func TestPlayerRedCard(t *testing.T) {
	player := engine.NewPlayer()
	player.SetSkills(*bc.Contracts, SkillsFromString(t, "40439920000726868070503716865792521545121682176182486071370780491777"))
	golden.Assert(t, dump.Sdump(player), t.Name()+".golden")
}
