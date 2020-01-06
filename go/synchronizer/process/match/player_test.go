package match_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
	"gotest.tools/assert"
)

func TestDefenceOfNullPlayer(t *testing.T) {
	p := match.NewNullPlayer()
	defence, err := p.Defence(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, defence, uint64(0))
}

func TestDefenceOfPlayer(t *testing.T) {
	p := match.NewPlayer("14606253788909032162646379450304996475079674564248175")
	defence, err := p.Defence(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, defence, uint64(955))
}
