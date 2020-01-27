package engine_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/engine"
)

var funcs = []struct {
	name string
	f    func(*testing.B, engine.Matches)
}{
	{"main", benchmarkPlay1stHalf},
	{"goRoutines", benchmarkPlay1stHalfParallel},
}

func BenchmarkPlayer1stHalf10(b *testing.B) {
	for _, f := range funcs {
		for i := 10; i < 200; i *= 2 {
			var matches engine.Matches
			for n := 0; n < i; n++ {
				matches = append(matches, *engine.NewMatch())
			}
			b.Run(fmt.Sprintf("%s/%d", f.name, i), func(b *testing.B) {
				f.f(b, matches)
			})
		}
	}
}

func benchmarkPlay1stHalf(b *testing.B, matches engine.Matches) {
	for n := 0; n < b.N; n++ {
		if err := matches.Play1stHalf(*bc.Contracts); err != nil {
			b.Error(err)
		}
	}
}

func benchmarkPlay1stHalfParallel(b *testing.B, matches engine.Matches) {
	for n := 0; n < b.N; n++ {
		if err := matches.Play1stHalfParallel(context.Background(), *bc.Contracts); err != nil {
			b.Error(err)
		}
	}
}
