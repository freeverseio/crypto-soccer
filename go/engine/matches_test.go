package engine_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/engine"
	"gotest.tools/assert"
)

var funcs = []struct {
	name string
	f    func(engine.Matches, context.Context, contracts.Contracts) error
}{
	{"main", engine.Matches.Play1stHalf},
	{"goRoutines", engine.Matches.Play1stHalfParallel},
}

func BenchmarkPlayer1stHalf(b *testing.B) {
	for _, f := range funcs {
		for i := 10; i < 500; i *= 2 {
			var matches engine.Matches
			for n := 0; n < i; n++ {
				matches = append(matches, *engine.NewMatch())
			}
			b.Run(fmt.Sprintf("%s/%d", f.name, i), func(b *testing.B) {
				f.f(matches, context.Background(), *bc.Contracts)
			})
		}
	}
}

func TestMatchesPlay1stHalfParallel(t *testing.T) {
	const nMatches = 100
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
