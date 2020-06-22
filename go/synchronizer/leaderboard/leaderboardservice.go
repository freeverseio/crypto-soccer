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

func ComputeLeague(
	contracts contracts.Contracts,
	matchDay int,
	matches []storage.Match,
) error {
	var teamIds [8]*big.Int
	var results [56][2]uint8

	bcLeaderboard, err := contracts.Leagues.ComputeLeagueLeaderBoard(
		&bind.CallOpts{},
		teamIds,
		results,
		uint8(matchDay),
	)
	if err != nil {
		return err
	}

	l := Leaderboard{}
	for i := range l {
		l[i].TeamId = teamIds[i].String()
		l[i].Position = int(bcLeaderboard.Ranking[i])
		l[i].Points = int(bcLeaderboard.Points[i].Int64())
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
		ComputeLeague(
			contracts,
			matchDay,
			matches[i:i+56],
		)
	}

	return nil
}
