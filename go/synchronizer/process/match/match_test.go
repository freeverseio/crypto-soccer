package match_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/match"
	"gotest.tools/assert"
)

func TestDefaultValues(t *testing.T) {
	t.Parallel()
	match, err := match.NewMatch(bc.Contracts)
	assert.NilError(t, err)
	assert.Assert(t, match != nil, "match is nil")
}

func TestPlay1stHalfWithDefaultValues(t *testing.T) {
	t.Parallel()
	match, _ := match.NewMatch(bc.Contracts)
	err := match.Play1stHalf()
	assert.NilError(t, err)
	assert.Equal(t, match.HomeGoals, uint8(0))
	assert.Equal(t, match.VisitorGoals, uint8(0))
	assert.Equal(t, match.HomeMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
	assert.Equal(t, match.VisitorMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
}

// TODO: reactive this test
// func TestPlay2ndHalfWithDefaultValues(t *testing.T) {
// 	match, _ := match.NewMatch(bc.Contracts)
// 	err := match.Play2ndHalf()
// 	assert.NilError(t, err)
// 	assert.Equal(t, match.HomeGoals, uint8(0))
// 	assert.Equal(t, match.VisitorGoals, uint8(0))
// 	assert.Equal(t, match.HomeMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
// 	assert.Equal(t, match.VisitorMatchLog.String(), "1645504557321206042155578968558872826709262232930097591983538176")
// }

func TestPlayi1stHalf(t *testing.T) {
	t.Parallel()
	m, _ := match.NewMatch(bc.Contracts)
	homePlayer := match.NewPlayerFromSkills("60912465658141224081372268432703414642709456376891023")
	visitorPlayer := match.NewPlayerFromSkills("527990852960211435545446683633031307934132992821212439")
	m.HomeTeam.Players[0] = homePlayer
	m.VisitorTeam.Players[0] = visitorPlayer
	err := m.Play1stHalf()
	assert.NilError(t, err)
	assert.Equal(t, m.HomeGoals, uint8(0))
	assert.Equal(t, m.VisitorGoals, uint8(0))
	assert.Equal(t, m.HomeMatchLog.String(), "68582984444590546630976961169593813219497174670109271642235310440448")
	assert.Equal(t, m.VisitorMatchLog.String(), "594466494909760143231211294687139552942416193784081023125303800627200")
	assert.Equal(t, m.HomeTeam.Players[0].Skills().String(), "60912471367131994905211792665847292440690001907877519")
	assert.Equal(t, m.HomeTeam.Players[1].Skills().String(), "0")
	assert.Equal(t, m.VisitorTeam.Players[0].Skills().String(), "527990858669202206369286207866175185732113538352198935")
	assert.Equal(t, m.VisitorTeam.Players[1].Skills().String(), "0")
}

func TestPlayi1stHalf_part2(t *testing.T) {
	t.Parallel()
	m, _ := match.NewMatch(bc.Contracts)
	m.Seed = [32]byte{0x1, 0x1f}
	m.StartTime = big.NewInt(34525345)
	m.HomeTeam.TeamID = big.NewInt(1)
	m.VisitorTeam.TeamID = big.NewInt(2)
	for i := 0; i < 11; i++ {
		m.HomeTeam.Players[i] = match.CreateDummyPlayer(t, bc.Contracts, 10, 10, 10, 10, 10)
		m.VisitorTeam.Players[i] = match.CreateDummyPlayer(t, bc.Contracts, 50, 50, 50, 50, 50)
	}
	err := m.Play1stHalf()
	assert.NilError(t, err)
	assert.Equal(t, m.HomeGoals, uint8(0))
	assert.Equal(t, m.VisitorGoals, uint8(0))
}
