package names_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/names"
)

func TestGeneratePlayerName(t *testing.T) {
	generator, err := names.New()
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	for i := 0; i < 10; i++ {
		name, err := generator.GeneratePlayerName(big.NewInt(int64(i)), big.NewInt(int64(5+i)))
		if err != nil {
			t.Fatalf("error generating name for player: %s", err)
		}
		fmt.Println(name)
		if len(name) == 0 {
			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
		}
	}
}
