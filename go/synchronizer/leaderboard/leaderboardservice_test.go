package leaderboard_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/storage/mock"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/leaderboard"
	"gotest.tools/assert"
)

func TestLeaderboardServiceSort(t *testing.T) {
	t.Run("empty matches", func(t *testing.T) {
		matches := []storage.Match{}
		leaderboard.Sort(matches)
		assert.Equal(t, len(matches), 0)
	})
	t.Run("one match", func(t *testing.T) {
		matches := []storage.Match{}
		matches = append(matches, storage.Match{})
		leaderboard.Sort(matches)
		assert.Equal(t, len(matches), 1)
		assert.Equal(t, matches[0], storage.Match{})
	})
	t.Run("complicate one", func(t *testing.T) {
		matches := []storage.Match{}
		match := storage.NewMatch()
		match.TimezoneIdx = 2
		matches = append(matches, *match)
		match.TimezoneIdx = 1
		match.CountryIdx = 1
		matches = append(matches, *match)
		match.CountryIdx = 0
		matches = append(matches, *match)
		leaderboard.Sort(matches)
		assert.Equal(t, matches[0].TimezoneIdx, uint8(1))
		assert.Equal(t, matches[0].CountryIdx, uint32(0))
		assert.Equal(t, matches[1].TimezoneIdx, uint8(1))
		assert.Equal(t, matches[1].CountryIdx, uint32(1))
		assert.Equal(t, matches[2].TimezoneIdx, uint8(2))

	})
}

func TestLeaderboardServiceNoMatches(t *testing.T) {
	timezone := 10
	matchDay := 0
	sto := mock.NewStorageService()
	sto.MatchStorageService.MatchesByTimezoneFunc = func(timezone uint8) ([]storage.Match, error) {
		return []storage.Match{}, nil
	}
	service := leaderboard.NewLeaderboardService(*sto)
	assert.NilError(t, service.Update(*bc.Contracts, matchDay, timezone))
}

func TestLeaderboardService1Match(t *testing.T) {
	timezone := 10
	matchDay := 0
	sto := mock.NewStorageService()
	sto.MatchStorageService.MatchesByTimezoneFunc = func(timezone uint8) ([]storage.Match, error) {
		return []storage.Match{storage.Match{}}, nil
	}
	service := leaderboard.NewLeaderboardService(sto)
	assert.Error(t, service.Update(*bc.Contracts, matchDay, timezone), "matches count not multiple 56")
}

func TestLeaderboardServiceLeague(t *testing.T) {
	timezone := 10
	matchDay := 0
	sto := mock.NewStorageService()
	sto.MatchStorageService.MatchesByTimezoneFunc = func(timezone uint8) ([]storage.Match, error) {
		return []storage.Match{storage.Match{}}, nil
	}
	service := leaderboard.NewLeaderboardService(sto)
	assert.Error(t, service.Update(*bc.Contracts, matchDay, timezone), "matches count not multiple 56")
}
