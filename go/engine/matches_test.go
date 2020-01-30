package engine_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/engine"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func BenchmarkPlayer1stHalfParallel(b *testing.B) {
	matchesCount := []int{10, 100, 200, 500, 1000, 2000}
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

func TestMatchesPlay1stHalf(t *testing.T) {
	const nMatches = 10
	var matches engine.Matches
	for i := 0; i < nMatches; i++ {
		matches = append(matches, *engine.NewMatch())
	}
	err := matches.Play1stHalf(context.Background(), *bc.Contracts)
	assert.NilError(t, err)
	golden.Assert(t, matches.DumpState(), t.Name()+".golden")
}

func TestMatchesPlay1stHalfParallel(t *testing.T) {
	t.Parallel()
	const nMatches = 10
	var matches engine.Matches
	for i := 0; i < nMatches; i++ {
		matches = append(matches, *engine.NewMatch())
	}
	err := matches.Play1stHalf(context.Background(), *bc.Contracts)
	assert.NilError(t, err)
	var matchesP engine.Matches
	for i := 0; i < nMatches; i++ {
		matchesP = append(matchesP, *engine.NewMatch())
	}
	err = matchesP.Play1stHalfParallel(context.Background(), *bc.Contracts)
	assert.Equal(t, matches.DumpState(), matchesP.DumpState())
}
