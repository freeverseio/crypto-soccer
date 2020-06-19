package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/synchronizer/leaderboard/leaderboardtest"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/leaderboard/postgres"
	"gotest.tools/assert"
)

func TestLeaderboardServiceInterface(t *testing.T) {
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	service := postgres.NewLeaderboardService(tx)
	leaderboardtest.TestLeaderboardServiceInterface(t, service)
}
