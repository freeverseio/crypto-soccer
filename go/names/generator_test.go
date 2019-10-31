package names_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/names"
)

func TestSubmitActionRoot(t *testing.T) {
	for i := 0; i < 10; i++ {
		name := names.GeneratePlayerName(big.NewInt(int64(i)))
		fmt.Println(name)
		if len(name) == 0 {
			t.Fatalf("Expecting non empty player name, but got \"%v\"", name)
		}
	}
}
