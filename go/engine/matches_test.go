package engine_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/engine"
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
		for i := 10; i < 200; i *= 2 {
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
