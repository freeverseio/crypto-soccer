package leaderboard_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/storage/memory"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/leaderboard"
	"gotest.tools/assert"
)

func TestLeaderboardServiceNoMatches(t *testing.T) {
	timezone := 10
	sto := memory.NewStorageService()
	service := leaderboard.NewLeaderboardService(sto)
	assert.NilError(t, service.Update(*bc.Contracts, timezone))
}

func TestLeaderboardService1Match(t *testing.T) {
	timezone := 10
	sto := memory.NewStorageService()
	sto.MatchService.Insert(storage.Match{})
	service := leaderboard.NewLeaderboardService(sto)
	assert.Error(t, service.Update(*bc.Contracts, timezone), "odd number of teams")
}
