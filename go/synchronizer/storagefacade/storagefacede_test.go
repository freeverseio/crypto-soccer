package storagefacade_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/storagefacade"
	"gotest.tools/assert"
)

func TestNoStorageMatches(t *testing.T) {
	tx, err := s.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	matches, err := storagefacade.NewMatchesByLeague(tx, 1, 0, 0, 0)
	assert.NilError(t, err)
	assert.Equal(t, len(matches), 0)
}
