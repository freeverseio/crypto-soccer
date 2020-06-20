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
	service *storage.StorageService
}

func NewLeaderboardService(service *storage.StorageService) *LeaderboardService {
	return &LeaderboardService{
		service: service,
	}
}

func (b LeaderboardService) Update(contracts contracts.Contracts, timezone int) error {
	matchDay := 0
	matches, err := b.service.MatchService.MatchesByTimezone(uint8(timezone))
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		return nil
	}

	if len(matches)%2 != 0 {
		return errors.New("matches count not multiple 56")
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].MatchDayIdx != matches[j].MatchDayIdx {
			return matches[i].MatchDayIdx > matches[j].MatchDayIdx
		}
		return matches[i].MatchIdx > matches[j].MatchIdx
	})

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
