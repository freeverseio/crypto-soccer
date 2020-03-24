package storage_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestTeamHistoryInsert(t *testing.T) {
	t.Parallel()
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	team, err := storage.TeamByTeamId(tx, teamID)
	assert.NilError(t, err)
	th := storage.NewTeamHistory(0, team)
	assert.NilError(t, th.Insert(tx))
	assert.Error(t, th.Insert(tx), "pq: duplicate key value violates unique constraint \"teams_histories_pkey\"")
}
