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

	matches = nil
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
	err = matches.Play1stHalfParallel(context.Background(), *bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, matches.DumpState(), t.Name()+".half.golden")
	err = matches.Play2ndHalfParallel(context.Background(), *bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, matches.DumpState(), t.Name()+".end.golden")
}
