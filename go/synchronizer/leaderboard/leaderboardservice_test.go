package leaderboard_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/storage/mock"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/leaderboard"
	"gotest.tools/assert"
)

func TestLeaderboardServiceNoMatches(t *testing.T) {
	timezone := 10
	sto := mock.NewStorageService()
	sto.MatchStorageService.MatchesByTimezoneFunc = func(timezone uint8) ([]storage.Match, error) {
		return []storage.Match{}, nil
	}
	service := leaderboard.NewLeaderboardService(*sto)
	assert.NilError(t, service.Update(*bc.Contracts, timezone))
}

// func TestLeaderboardService1Match(t *testing.T) {
// 	timezone := 10
// 	sto := mock.NewStorageService()
// 	service := leaderboard.NewLeaderboardService(sto)
// 	assert.Error(t, service.Update(*bc.Contracts, timezone), "odd number of teams")
// }
