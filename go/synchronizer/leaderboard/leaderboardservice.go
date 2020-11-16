package leaderboard

import (
	"math/big"
	"sort"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-log.v1"

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
			return m0.MatchDayIdx < m1.MatchDayIdx
		}
		return m0.MatchIdx < m1.MatchIdx
	})
}

func UpdateLeagueLeaderboard(
	contracts contracts.Contracts,
	matchDay int,
	matches [56]storage.Match,
	teams [8]storage.Team,
) ([8]storage.Team, error) {
	// log.Infof("UpdateLeagueLeaderboard matches %+v, teams %+v", matches, teams)

	timezoneIdx := matches[0].TimezoneIdx
	countryIdx := matches[0].CountryIdx
	leagueIdx := matches[0].LeagueIdx

	for _, match := range matches {
		if match.TimezoneIdx != timezoneIdx {
			return [8]storage.Team{}, errors.New("matches of different timezone")
		}
		if match.CountryIdx != countryIdx {
			return [8]storage.Team{}, errors.New("matches of different country")
		}
		if match.LeagueIdx != leagueIdx {
			return [8]storage.Team{}, errors.New("matches of different league")
		}
	}

	for i := range teams {
		if teams[i].TeamIdxInLeague != uint32(i) {
			return [8]storage.Team{}, errors.New("not ordered team")
		}
	}

	var results [56][2]uint8
	for i := range matches {
		results[i][0] = matches[i].HomeGoals
		results[i][1] = matches[i].VisitorGoals
	}
	var teamIdxInLeague [8]*big.Int
	for i := range teamIdxInLeague {
		teamIdxInLeague[i] = big.NewInt(int64(i))
	}

	// log.Infof("Calling ComputeLeagueLeaderboard %v %v %v", teamIdxInLeague, results, matchDay)
	llb, err := contracts.Leagues.ComputeLeagueLeaderBoard(
		&bind.CallOpts{},
		teamIdxInLeague,
		results,
		uint8(matchDay),
	)
	if err != nil {
		return [8]storage.Team{}, errors.Wrapf(err, "failed calling the BC teamIdxInLeague %v , results %v, matchDay %v", teamIdxInLeague, results, matchDay)
	}

	for i := 0; i < 8; i++ {
		teams[i].LeaderboardPosition = int(llb.Ranking[i])
	}

	return teams, nil
}

func (b LeaderboardService) UpdateTimezoneLeaderboards(
	contracts contracts.Contracts,
	timezone int,
	matchDay int,
) error {
	log.Debugf("UpdateTimezoneLeaderboard timezone %v matchDay %v", timezone, matchDay)
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
		leagueMatches := [56]storage.Match{}
		copy(leagueMatches[:], matches[i:i+56])

		timezoneIdx := leagueMatches[0].TimezoneIdx
		countryIdx := leagueMatches[0].CountryIdx
		leagueIdx := leagueMatches[0].LeagueIdx
		teams, err := b.service.TeamService().TeamsByTimezoneIdxCountryIdxLeagueIdx(timezoneIdx, countryIdx, leagueIdx)
		if err != nil {
			return err
		}
		if len(teams) != 8 {
			return errors.New("number of teams of a league has to be 8")
		}
		// ordering by index in league
		sort.Slice(teams[:], func(i, j int) bool {
			return teams[i].TeamIdxInLeague < teams[j].TeamIdxInLeague
		})
		leagueTeams := [8]storage.Team{}
		copy(leagueTeams[:], teams)
		if leagueTeams, err = UpdateLeagueLeaderboard(
			contracts,
			matchDay,
			leagueMatches,
			leagueTeams,
		); err != nil {
			return errors.Wrapf(err, "failed update leaderboard timezone %v, country %v, league %v", timezoneIdx, countryIdx, leagueIdx)
		}

		for i := range leagueTeams {
			if err := b.service.TeamService().UpdateLeaderboardPosition(leagueTeams[i].TeamID, leagueTeams[i].LeaderboardPosition); err != nil {
				return err
			}
		}
	}

	return nil
}
