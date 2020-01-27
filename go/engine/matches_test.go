package engine_test

import (
	"context"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/engine"
)

func BenchmarkPlayer1stHalf10(b *testing.B) {
	benchmarkPlay1stHalf(b, 10)
}

func BenchmarkPlayer1stHalf20(b *testing.B) {
	benchmarkPlay1stHalf(b, 20)
}

func BenchmarkPlayer1stHalf40(b *testing.B) {
	benchmarkPlay1stHalf(b, 40)
}

func BenchmarkPlayer1stHalf80(b *testing.B) {
	benchmarkPlay1stHalf(b, 80)
}

func BenchmarkPlayer1stHalf160(b *testing.B) {
	benchmarkPlay1stHalf(b, 160)
}

func BenchmarkPlayer1stHalfParallel10(b *testing.B) {
	benchmarkPlay1stHalfParallel(b, 10)
}

func BenchmarkPlayer1stHalfParallel20(b *testing.B) {
	benchmarkPlay1stHalfParallel(b, 20)
}

func BenchmarkPlayer1stHalfParallel40(b *testing.B) {
	benchmarkPlay1stHalfParallel(b, 40)
}

func BenchmarkPlayer1stHalfParallel80(b *testing.B) {
	benchmarkPlay1stHalfParallel(b, 80)
}

func BenchmarkPlayer1stHalfParallel160(b *testing.B) {
	benchmarkPlay1stHalfParallel(b, 160)
}

func BenchmarkPlayer1stHalfParallel320(b *testing.B) {
	benchmarkPlay1stHalfParallel(b, 320)
}

func benchmarkPlay1stHalf(b *testing.B, nMatches int) {
	var matches engine.Matches
	for i := 0; i < nMatches; i++ {
		matches = append(matches, *engine.NewMatch())
	}
	for n := 0; n < b.N; n++ {
		if err := matches.Play1stHalf(*bc.Contracts); err != nil {
			b.Error(err)
		}
	}
}

func benchmarkPlay1stHalfParallel(b *testing.B, nMatches int) {
	var matches engine.Matches
	for i := 0; i < nMatches; i++ {
		matches = append(matches, *engine.NewMatch())
	}
	for n := 0; n < b.N; n++ {
		if err := matches.Play1stHalfParallel(context.Background(), *bc.Contracts); err != nil {
			b.Error(err)
		}
	}
}
