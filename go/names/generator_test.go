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
	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
	var timezone uint8
	var countryIdxInTZ uint64
	for i := 0; i < 10; i++ {
		playerId := big.NewInt(int64(i))
		timezone = 19
		countryIdxInTZ = 0
		name, err := generator.GeneratePlayerFullName(playerId, timezone, countryIdxInTZ)
		if err != nil {
			t.Fatalf("error generating name for player: %s", err)
		}
		fmt.Println(name)
		if len(name) == 0 {
			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
		}
	}
}
