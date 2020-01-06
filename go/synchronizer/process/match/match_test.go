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
	homePlayer := match.NewPlayer("60912465658141224081372268432703414642709456376891023")
	visitorPlayer := match.NewPlayer("527990852960211435545446683633031307934132992821212439")
	m.HomeTeam.Players[0] = homePlayer
	m.VisitorTeam.Players[0] = visitorPlayer
	// err := m.Play1stHalf()
	_, err := m.Process(false)
	assert.NilError(t, err)
	assert.Equal(t, m.HomeGoals, uint8(0))
	assert.Equal(t, m.VisitorGoals, uint8(0))
	assert.Equal(t, m.HomeMatchLog.String(), "68582984444590546630976961169593813219497174670109271642235310440448")
	assert.Equal(t, m.VisitorMatchLog.String(), "68582984444590546630976961169593813219497174670109271642235310440448")
	assert.Equal(t, m.HomeTeam.Players[0].Skills().String(), "60912471367131994905211792665847292440690001907877519")
	assert.Equal(t, m.HomeTeam.Players[1].Skills().String(), "0")
	assert.Equal(t, m.VisitorTeam.Players[0].Skills().String(), "60912471367131994905211792665847292440690001907877519")
	assert.Equal(t, m.VisitorTeam.Players[1].Skills().String(), "0")
}
