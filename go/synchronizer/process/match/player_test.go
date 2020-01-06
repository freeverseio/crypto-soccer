package match_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
	"gotest.tools/assert"
)

func TestSkillsOfNullPlayer(t *testing.T) {
	p := match.NewNullPlayer()
	defence, err := p.Defence(bc.Contracts.Assets)
	assert.NilError(t, err)
	assert.Equal(t, defence, uint64(0))
}
