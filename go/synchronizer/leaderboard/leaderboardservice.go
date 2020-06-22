package leaderboard

import (
	"errors"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

type LeaderboardService struct {
	service storage.StorageService
}

func NewLeaderboardService(service storage.StorageService) *LeaderboardService {
	return &LeaderboardService{
		service: service,
	}
}

func Sort(matches []storage.Match) {
	sort.Slice(matches, func(i, j int) bool {
		m0 := matches[i]
		m1 := matches[j]
		if m0.TimezoneIdx != m1.TimezoneIdx {
			return m0.TimezoneIdx < m1.TimezoneIdx
		}
		if m0.CountryIdx != m1.CountryIdx {
			return m0.CountryIdx < m1.CountryIdx
		}
		if m0.LeagueIdx != m1.LeagueIdx {
			return m0.LeagueIdx < m1.LeagueIdx
		}
		if m0.MatchDayIdx != m1.MatchDayIdx {
			return m0.MatchDayIdx < m0.MatchDayIdx
		}
		return m0.MatchIdx < m0.MatchIdx
	})
}

func (b LeaderboardService) ComputeLeague(
	contracts contracts.Contracts,
	matchDay int,
	matches []storage.Match,
) error {
	if len(matches) != 56 {
		return errors.New("number of matches in not 56")
	}

	var results [56][2]uint8
	for i := range matches {
		results[i][0] = matches[i].HomeGoals
		results[i][1] = matches[i].VisitorGoals
	}
	var teamIds [8]*big.Int

	llb, err := contracts.Leagues.ComputeLeagueLeaderBoard(
		&bind.CallOpts{},
		teamIds,
		results,
		uint8(matchDay),
	)
	if err != nil {
		return err
	}

	for i := 0; i < 8; i++ {
		if err := b.service.TeamService().UpdateLeaderboardPosition(teamIds[i].String(), int(llb.Ranking[i])); err != nil {
			return err
		}
	}

	return nil
}

func (b LeaderboardService) Update(
	contracts contracts.Contracts,
	timezone int,
	matchDay int,
) error {
	matches, err := b.service.MatchService().MatchesByTimezone(uint8(timezone))
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		return nil
	}

	if len(matches)%56 != 0 {
		return errors.New("matches count not multiple 56")
	}

	Sort(matches)

	for i := 0; i < len(matches); i += 56 {
		b.ComputeLeague(
			contracts,
			matchDay,
			matches[i:i+56],
		)
	}

	return nil
}
