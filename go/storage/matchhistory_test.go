package storage_test

import (
	"math/big"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestMatchHistory(t *testing.T) {
	t.Parallel()
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	blockNumber := uint64(10)
	match := storage.NewMatch()
	match.TimezoneIdx = timezoneIdx
	match.HomeTeamID, _ = new(big.Int).SetString(teamID, 10)
	match.VisitorTeamID, _ = new(big.Int).SetString(teamID1, 10)
	match.StartEpoch = 5
	mh := storage.NewMatchHistory(blockNumber, *match)

	assert.NilError(t, mh.Insert(tx))
}
