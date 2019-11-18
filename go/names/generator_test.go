package names_test

import (
	"fmt"
	"hash/fnv"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/names"
)

func int_hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func TestGeneratePlayerName(t *testing.T) {
	generator, err := names.New("./sql/00_goalRev.db")
	if err != nil {
		t.Fatalf("error creating database for player names: %s", err)
	}
	// WARNING: both timezone and countryIdxInTZ are derivable from playerId
	var timezone uint8
	var countryIdxInTZ uint64
	var result string = ""
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
		result += name
	}
	if int_hash(result) != uint64(652704827421104006) {
		fmt.Println("the just-obtained hash is: ")
		fmt.Println(int_hash(result))
		t.Fatal("result of generating names not as expected")
	}
}
