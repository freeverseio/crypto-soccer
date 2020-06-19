package leaderboard_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage/memory"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/leaderboard"
	"gotest.tools/assert"
)

func TestLeaderboardService(t *testing.T) {
	sService := memory.NewStorageService()
	lService := leaderboard.NewLeaderboardService(sService)
	timezone := 10
	matchDay := 3
	_, err := lService.Compute(*bc.Contracts, timezone, matchDay)
	assert.Error(t, err, "")
}
