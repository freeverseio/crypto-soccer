package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestTacticHistoryInsertWithSameBlockNumber(t *testing.T) {
	t.Parallel()
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	tactic := storage.NewTactic()
	th := storage.NewTacticHistory(0, *tactic)
	th.TeamID = teamID
	assert.NilError(t, th.Insert(tx))
	assert.Error(t, th.Insert(tx), "pq: duplicate key value violates unique constraint \"tactics_histories_pkey\"")
}

func TestTacticHistoryInsert(t *testing.T) {
	t.Parallel()
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	tactic := storage.NewTactic()
	th := storage.NewTacticHistory(0, *tactic)
	th.TeamID = teamID
	assert.NilError(t, th.Insert(tx))
	th.BlockNumber = 1
	assert.NilError(t, th.Insert(tx))
}
