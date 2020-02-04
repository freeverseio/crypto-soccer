package engine_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/engine"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func BenchmarkPlayer1stHalfParallel(b *testing.B) {
	matchesCount := []int{50, 100, 200, 400, 800, 1600, 3200}
	for _, count := range matchesCount {
		b.Run(fmt.Sprintf("%d", count), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				var matches engine.Matches
				for n := 0; n < count; n++ {
					matches = append(matches, *engine.NewMatch())
				}
				b.StartTimer()
				matches.Play1stHalfParallel(context.Background(), *bc.Contracts)
			}
		})
	}
}

// func TestMatchesPlay1stHalf(t *testing.T) {
// 	const nMatches = 10
// 	var matches engine.Matches
// 	for i := 0; i < nMatches; i++ {
// 		matches = append(matches, *engine.NewMatch())
// 	}
// 	err := matches.Play1stHalf(*bc.Contracts)
// 	assert.NilError(t, err)
// 	golden.Assert(t, matches.DumpState(), t.Name()+".golden")
// }

// func TestMatchesPlay1stHalfParallel(t *testing.T) {
// 	t.Parallel()
// 	const nMatches = 10
// 	var matches engine.Matches
// 	for i := 0; i < nMatches; i++ {
// 		matches = append(matches, *engine.NewMatch())
// 	}
// 	err := matches.Play1stHalf(*bc.Contracts)
// 	assert.NilError(t, err)
// 	var matchesP engine.Matches
// 	for i := 0; i < nMatches; i++ {
// 		matchesP = append(matchesP, *engine.NewMatch())
// 	}
// 	err = matchesP.Play1stHalfParallel(context.Background(), *bc.Contracts)
// 	assert.Equal(t, matches.DumpState(), matchesP.DumpState())
// }

func TestMatchesPlaySequentialAndPlayParallal(t *testing.T) {
	t.Parallel()
	var matches engine.Matches
	for i := 0; i < 2; i++ {
		match := engine.NewMatch()
		match.Seed = sha256.Sum256([]byte(fmt.Sprintf("%d", i)))
		for i := 0; i < 25; i++ {
			var err error
			match.HomeTeam.Players[i], err = engine.NewPlayerFromSkills(*bc.Contracts, "16573429227295117480385309339445376240739796176995438")
			assert.NilError(t, err)
			match.VisitorTeam.Players[i], err = engine.NewPlayerFromSkills(*bc.Contracts, "16573429227295117480385309340654302060354425351701614")
			assert.NilError(t, err)
		}
		matches = append(matches, *match)
	}
	golden.Assert(t, matches.DumpState(), t.Name()+".begin.golden")
	err := matches.Play1stHalf(*bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, matches.DumpState(), t.Name()+".half.golden")
	err = matches.Play2ndHalf(*bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, matches.DumpState(), t.Name()+".end.golden")
}

// func TestNewMatchByLeagueWithNoMatches(t *testing.T) {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer tx.Rollback()
// 	timezoneIdx := uint8(1)
// 	day := uint8(0)
// 	matches, err := engine.NewMatchesFromTimezoneIdxMatchdayIdx(tx, timezoneIdx, day)
// 	assert.NilError(t, err)
// 	assert.Equal(t, len(*matches), 0)
// }

// func TestMatchesFromStorage(t *testing.T) {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer tx.Rollback()
// 	createMatches(t, tx)
// 	timezoneIdx := uint8(1)
// 	day := uint8(0)
// 	matches, err := engine.NewMatchesFromTimezoneIdxMatchdayIdx(tx, timezoneIdx, day)
// 	assert.NilError(t, err)
// 	assert.Equal(t, len(*matches), 8)
// 	match := (*matches)[0]
// 	assert.Equal(t, match.HomeTeam.TeamID.String(), "10")
// 	assert.Equal(t, match.VisitorTeam.TeamID.String(), "10")
// 	assert.Equal(t, match.HomeMatchLog.String(), "12")
// 	assert.Equal(t, match.VisitorMatchLog.String(), "13")
// 	assert.Equal(t, match.HomeGoals, uint8(0))
// 	assert.Equal(t, match.VisitorGoals, uint8(0))
// 	assert.Equal(t, match.HomeTeam.Players[0].Skills().String(), "123456")
// 	assert.Equal(t, match.VisitorTeam.Players[0].Skills().String(), "123456")
// 	assert.Equal(t, len(match.Events), 0)
// }

// func createMatches(t *testing.T, tx *sql.Tx) {
// 	timezoneIdx := uint8(1)
// 	timezone := storage.Timezone{timezoneIdx}
// 	err := timezone.Insert(tx)
// 	assert.NilError(t, err)

// 	countryIdx := uint32(0)
// 	country := storage.Country{timezone.TimezoneIdx, countryIdx}
// 	err = country.Insert(tx)
// 	assert.NilError(t, err)

// 	leagueIdx := uint32(0)
// 	league := storage.League{timezoneIdx, countryIdx, leagueIdx}
// 	err = league.Insert(tx)
// 	assert.NilError(t, err)

// 	var team storage.Team
// 	team.TeamID = big.NewInt(10)
// 	team.TimezoneIdx = timezoneIdx
// 	team.CountryIdx = countryIdx
// 	team.Owner = "ciao"
// 	team.LeagueIdx = leagueIdx
// 	err = team.Insert(tx)
// 	assert.NilError(t, err)

// 	var player storage.Player
// 	player.TeamId = team.TeamID
// 	player.EncodedSkills = big.NewInt(123456)
// 	assert.NilError(t, player.Insert(tx))

// 	for i := 0; i < 8; i++ {
// 		matchDayIdx := uint8(0)
// 		matchIdx := uint8(i)
// 		match := storage.Match{
// 			TimezoneIdx:     timezoneIdx,
// 			CountryIdx:      countryIdx,
// 			LeagueIdx:       leagueIdx,
// 			MatchDayIdx:     matchDayIdx,
// 			MatchIdx:        matchIdx,
// 			HomeTeamID:      big.NewInt(10),
// 			VisitorTeamID:   big.NewInt(10),
// 			HomeMatchLog:    big.NewInt(12),
// 			VisitorMatchLog: big.NewInt(13),
// 		}
// 		err = match.Insert(tx)
// 		assert.NilError(t, err)
// 	}
// }
