package leaderboard_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/storage/mock"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/leaderboard"
	"gotest.tools/assert"
	"gotest.tools/golden"
)

func TestLeaderboardServiceFrom1200NoMatches(t *testing.T) {
	timezone := 10
	matchDay := 0
	sto := mock.NewStorageService()
	sto.MatchStorageService.MatchesByTimezoneFunc = func(timezone uint8) ([]storage.Match, error) {
		return []storage.Match{}, nil
	}
	service := leaderboard.NewLeaderboardService(*sto)
	assert.NilError(t, service.UpdateTimezoneLeaderboardsFrom1200(*bc.Contracts, matchDay, timezone))
}

func TestLeaderboardServiceFrom12001Match(t *testing.T) {
	timezone := 10
	matchDay := 0
	sto := mock.NewStorageService()
	sto.MatchStorageService.MatchesByTimezoneFunc = func(timezone uint8) ([]storage.Match, error) {
		return []storage.Match{storage.Match{}}, nil
	}
	service := leaderboard.NewLeaderboardService(sto)
	assert.Error(t, service.UpdateTimezoneLeaderboardsFrom1200(*bc.Contracts, matchDay, timezone), "matches count not multiple 56")
}

func TestLeaderboardServiceFrom1200League(t *testing.T) {
	matches := [56]storage.Match{}
	for i := range matches {
		matches[i] = *storage.NewMatch()
	}
	teams := [8]storage.Team{}
	for i := range teams {
		teams[i] = *storage.NewTeam()
		teams[i].TeamID = fmt.Sprintf("%d", i)
		teams[i].TeamIdxInLeague = uint32(i)
	}

	sto := mock.NewStorageService()
	sto.MatchStorageService.MatchesByTimezoneFunc = func(timezone uint8) ([]storage.Match, error) {
		return matches[:], nil
	}
	sto.TeamStorageService.TeamsByTimezoneIdxCountryIdxLeagueIdxFunc = func(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) ([]storage.Team, error) {
		return teams[:], nil
	}
	sto.TeamStorageService.UpdateLeaderboardPositionFunc = func(teamId string, position int) error {
		id, err := strconv.Atoi(teamId)
		if err != nil {
			return err
		}
		teams[id].LeaderboardPosition = position
		return nil
	}

	service := leaderboard.NewLeaderboardService(sto)
	timezone := 10
	matchDay := 15
	assert.NilError(t, service.UpdateTimezoneLeaderboardsFrom1200(*bc.Contracts, matchDay, timezone))

	golden.Assert(t, dump.Sdump(teams), t.Name()+".golden")
}

func TestLeaderboardServiceUpdateLeagueLeaderboardFrom1200(t *testing.T) {
	matches := [56]storage.Match{}
	for i := range matches {
		matches[i] = *storage.NewMatch()
	}
	teams := [8]storage.Team{}
	for i := range teams {
		teams[i] = *storage.NewTeam()
		teams[i].TeamIdxInLeague = uint32(i)
	}
	t.Run("matchDay0AllDraw", func(t *testing.T) {
		matchDay := 0
		rTeams, err := leaderboard.UpdateLeagueLeaderboardFrom1200(
			*bc.Contracts,
			matchDay,
			matches,
			teams,
		)
		assert.NilError(t, err)
		golden.Assert(t, dump.Sdump(rTeams), t.Name()+".golden")
	})
	t.Run("matchDay0VisitorWins", func(t *testing.T) {
		matchDay := 0
		matches[0].VisitorGoals = 3
		rTeams, err := leaderboard.UpdateLeagueLeaderboardFrom1200(
			*bc.Contracts,
			matchDay,
			matches,
			teams,
		)
		assert.NilError(t, err)
		golden.Assert(t, dump.Sdump(rTeams), t.Name()+".golden")
	})
}
