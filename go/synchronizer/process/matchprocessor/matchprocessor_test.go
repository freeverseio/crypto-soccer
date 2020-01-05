package matchprocessor_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/process/matchprocessor"
)

func TestDefaultValues(t *testing.T) {
	if mp := matchprocessor.NewMatch(bc.Contracts); mp == nil {
		t.Fatal("New instance is nil")
	}
}
