package process_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/engine"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func BenchmarkPlayer1stHalfParallel(b *testing.B) {
	matchesCount := []int{50, 100, 200, 400, 800, 1600, 3200}
	for _, count := range matchesCount {
		b.Run(fmt.Sprintf("%d", count), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				var matches process.Matches
				for n := 0; n < count; n++ {
					matches = append(matches, *engine.NewMatch())
				}
				b.StartTimer()
				matches.Play1stHalfParallel(context.Background(), *bc.Contracts)
			}
		})
	}
}

func TestMatchesPlaySequentialAndPlayParallal(t *testing.T) {
	t.Parallel()
	var matches process.Matches
	for i := 0; i < 2; i++ {
		match := engine.NewMatch()
		match.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
		match.Seed = sha256.Sum256([]byte(fmt.Sprintf("%d", i)))
		for i := 0; i < 25; i++ {
			match.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309339445376240739796176995438"))
			match.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309340654302060354425351701614"))
		}
		matches = append(matches, *match)
	}
	golden.Assert(t, dump.Sdump(matches), t.Name()+".begin.golden")
	for i := 0; i < len(matches); i++ {
		assert.NilError(t, matches[i].Play1stHalf(*bc.Contracts))
	}
	golden.Assert(t, dump.Sdump(matches), t.Name()+".half.golden")
	for i := 0; i < len(matches); i++ {
		assert.NilError(t, matches[i].Play2ndHalf(*bc.Contracts))
	}
	golden.Assert(t, dump.Sdump(matches), t.Name()+".end.golden")

	matches = nil
	for i := 0; i < 2; i++ {
		match := engine.NewMatch()
		match.StartTime = big.NewInt(1570147200 + 3600*24*365*7)
		match.Seed = sha256.Sum256([]byte(fmt.Sprintf("%d", i)))
		for i := 0; i < 25; i++ {
			match.HomeTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309339445376240739796176995438"))
			match.VisitorTeam.Players[i].SetSkills(*bc.Contracts, SkillsFromString(t, "16573429227295117480385309340654302060354425351701614"))
		}
		matches = append(matches, *match)
	}
	golden.Assert(t, dump.Sdump(matches), t.Name()+".begin.golden")
	assert.NilError(t, matches.Play1stHalfParallel(context.Background(), *bc.Contracts))
	golden.Assert(t, dump.Sdump(matches), t.Name()+".half.golden")
	assert.NilError(t, matches.Play2ndHalfParallel(context.Background(), *bc.Contracts))
	golden.Assert(t, dump.Sdump(matches), t.Name()+".end.golden")
}

func TestMatchesSetTactics(t *testing.T) {
	t.Parallel()
	var matches process.Matches
	teamID := "3"
	matches = append(matches, *engine.NewMatch())
	matches = append(matches, *engine.NewMatch())
	matches[0].HomeTeam.TeamID = teamID
	matches[1].VisitorTeam.TeamID = teamID
	golden.Assert(t, dump.Sdump(matches), t.Name()+".begin.golden")
	tactics := []storage.Tactic{}
	tactic := storage.Tactic{teamID, 1, 0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 25, 11, 0, 25, 11, 0, 25, 11, 0, false, false, true, false, false, true, false, false, false, false}
	tactics = append(tactics, tactic)
	assert.NilError(t, matches.SetTactics(*bc.Contracts, tactics))
	golden.Assert(t, dump.Sdump(matches), t.Name()+".end.golden")
}

func TestMatchesSetTrainings(t *testing.T) {
	t.Parallel()
	var matches process.Matches
	teamID := "3"
	matches = append(matches, *engine.NewMatch())
	matches = append(matches, *engine.NewMatch())
	matches[0].HomeTeam.TeamID = teamID
	matches[1].VisitorTeam.TeamID = teamID
	golden.Assert(t, dump.Sdump(matches), t.Name()+".begin.golden")
	trainings := []storage.Training{}
	training := storage.Training{}
	training.TeamID = teamID
	trainings = append(trainings, training)
	assert.NilError(t, matches.SetTrainings(*bc.Contracts, trainings))
	golden.Assert(t, dump.Sdump(matches), t.Name()+".end.golden")
}

func TestMinute2Round(t *testing.T) {
	t.Parallel()
	cases := []struct {
		Minute int
		Round  uint8
	}{
		{0, 0},
		{1, 0},
		{4, 0},
		{5, 1},
		{9, 1},
		{10, 2},
		{14, 2},
		{15, 3},
		{19, 3},
		{20, 5},
		{44, 10},
		{45, 11},
		{46, 0},
		{49, 0},
		{50, 1},
		{89, 10},
		{90, 11},
		{91, 11},
		{100, 11},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			assert.Equal(t, process.Minute2Round(tc.Minute), tc.Round)
		})
	}
}

// func TestMatchesFromToStorage(t *testing.T) {
// 	t.Parallel()
// 	tx, err := universedb.Begin()
// 	assert.NilError(t, err)
// 	defer tx.Rollback()

// 	timezoneIdx := uint8(1)
// 	countryIdx := big.NewInt(0)
// 	divisionIdx := big.NewInt(0)
// 	day := uint8(0)
// 	divisionCreationProcessor, err := process.NewDivisionCreationProcessor(bc.Contracts, namesdb)
// 	assert.NilError(t, err)
// 	assert.NilError(t, divisionCreationProcessor.Process(tx, assets.AssetsDivisionCreation{timezoneIdx, countryIdx, divisionIdx, types.Log{}}))
// 	matches, err := process.NewMatchesFromTimezoneIdxMatchdayIdx(tx, timezoneIdx, day)
// 	match := (*matches)[0]
// 	match.Seed = [32]byte{0x4}
// 	assert.NilError(t, err)
// 	golden.Assert(t, dump.Sdump(match), t.Name()+".start.golden")
// 	assert.NilError(t, match.Play1stHalf(*bc.Contracts))
// 	golden.Assert(t, dump.Sdump(match), t.Name()+".half.golden")

// 	beginPlayer, err := storage.PlayerByPlayerId(tx, big.NewInt(274877906944))
// 	assert.NilError(t, err)
// 	assert.NilError(t, match.ToStorage(*bc.Contracts, tx))
// 	halfPlayer, err := storage.PlayerByPlayerId(tx, big.NewInt(274877906944))
// 	assert.NilError(t, err)
// 	assert.Assert(t, beginPlayer.Defence != halfPlayer.Defence)
// }
