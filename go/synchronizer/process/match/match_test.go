package match_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
	"gotest.tools/assert"
)

func TestDefaultValues(t *testing.T) {
	match, err := match.NewMatch(bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, match != nil, "match is nil")
}

func TestPlay1stHalfWithDefaultValues(t *testing.T) {
	match, _ := match.NewMatch(bc.Contracts)
	_, err := match.Process(false)
	assert.NilError(t, err)
	assert.Equal(t, match.HomeGoals, uint8(0))
	assert.Equal(t, match.VisitorGoals, uint8(0))
	assert.Equal(t, match.HomeMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
	assert.Equal(t, match.VisitorMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
}

func TestPlayi1stHalf(t *testing.T) {
	m, _ := match.NewMatch(bc.Contracts)
	for i := 0; i < 25; i++ {
		m.HomeTeam.Players[i] = match.NewPlayer("60912465658141224081372268432703414642709456376891023")
		m.VisitorTeam.Players[i] = match.NewPlayer("60912465658141224081372268432703414642709456376891023")
	}
	is2ndHalf := false
	_, err := m.Process(is2ndHalf)
	assert.NilError(t, err)
	assert.Equal(t, m.HomeGoals, uint8(0))
	assert.Equal(t, m.VisitorGoals, uint8(0))
	assert.Equal(t, m.HomeMatchLog.String(), "754396374849259078542209549811211635835627530328040412055968287817728")
	assert.Equal(t, m.VisitorMatchLog.String(), "754396374849259078542209549811211635835627530328040412055968287817728")
}
