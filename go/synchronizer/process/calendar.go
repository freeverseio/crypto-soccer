package process

import (
	"database/sql"
	"errors"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	log "github.com/sirupsen/logrus"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

type Calendar struct {
	contracts *contracts.Contracts
}

func NewCalendar(contracts *contracts.Contracts) *Calendar {
	return &Calendar{contracts}
}

func (b *Calendar) Generate(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) error {
	league, err := storage.LeagueByLeagueIdx(tx, leagueIdx)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Unexistent league")
	}

	for matchDay := uint8(0); matchDay < contracts.MatchDays; matchDay++ {
		for match := uint8(0); match < contracts.MatchesPerDay; match++ {
			m := storage.NewMatch()
			m.TimezoneIdx = timezoneIdx
			m.CountryIdx = countryIdx
			m.LeagueIdx = leagueIdx
			m.MatchDayIdx = matchDay
			m.MatchIdx = match
			err = m.Insert(tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b Calendar) GetAllMatchdaysUTCInCurrentRound(timezoneIdx uint8, verse *big.Int) ([14]*big.Int, error) {
	tz1, err := b.contracts.Updates.GetTimeZoneForRound1(&bind.CallOpts{})
	if err != nil {
		return [14]*big.Int{}, err
	}
	round, err := b.contracts.Updates.GetCurrentRoundPure(&bind.CallOpts{}, timezoneIdx, tz1, verse)
	if err != nil {
		return [14]*big.Int{}, err
	}
	matchesStart, err := b.contracts.Updates.GetAllMatchdaysUTCInRound(&bind.CallOpts{}, timezoneIdx, round)
	if err != nil {
		return [14]*big.Int{}, err
	}
	return matchesStart, nil
}

func (b Calendar) GetAllMatchdaysUTCInNextRound(timezoneIdx uint8, verse *big.Int) ([14]*big.Int, error) {
	tz1, err := b.contracts.Updates.GetTimeZoneForRound1(&bind.CallOpts{})
	if err != nil {
		return [14]*big.Int{}, err
	}
	round, err := b.contracts.Updates.GetCurrentRoundPure(&bind.CallOpts{}, timezoneIdx, tz1, verse)
	if err != nil {
		return [14]*big.Int{}, err
	}
	round.Add(round, big.NewInt(1))
	matchesStart, err := b.contracts.Updates.GetAllMatchdaysUTCInRound(&bind.CallOpts{}, timezoneIdx, round)
	if err != nil {
		return [14]*big.Int{}, err
	}
	return matchesStart, nil
}

func shiftBack(t uint8) uint8 {
	if t < TEAMS_PER_LEAGUE {
		return t
	}

	return t - (TEAMS_PER_LEAGUE - 1)
}

func getTeamsInMatchFirstHalf(matchday uint8, matchIdxInDay uint8) (uint8, uint8) {
	team1 := uint8(0)
	if matchIdxInDay > 0 {
		team1 = shiftBack(TEAMS_PER_LEAGUE - matchIdxInDay + matchday)
	}

	team2 := uint8(shiftBack(matchIdxInDay + 1 + matchday))
	if (matchday % 2) == 0 {
		return team1, team2
	}
	return team2, team1
}

type teamsIdxDuple struct {
	HomeIdx    uint8
	VisitorIdx uint8
}

func (b *Calendar) getTeamsInLeagueMatch(matchday uint8, matchIdxInDay uint8) (teamsIdxDuple, error) {
	MATCHDAYS := uint8(14)
	MATCHES_PER_DAY := uint8(4)
	homeIdx := uint8(0)
	visitorIdx := uint8(0)

	if matchday > MATCHDAYS {
		return teamsIdxDuple{homeIdx, visitorIdx}, errors.New("wrong match day")
	}

	if matchIdxInDay > MATCHES_PER_DAY {
		return teamsIdxDuple{homeIdx, visitorIdx}, errors.New("wrong match")
	}
	teamsDup := teamsIdxDuple{0, 0}
	if matchday < (TEAMS_PER_LEAGUE - 1) {
		homeIdx, visitorIdx := getTeamsInMatchFirstHalf(matchday, matchIdxInDay)
		teamsDup.HomeIdx = homeIdx
		teamsDup.VisitorIdx = visitorIdx
	} else {
		visitorIdx, homeIdx := getTeamsInMatchFirstHalf(matchday-(TEAMS_PER_LEAGUE-1), matchIdxInDay)
		teamsDup.HomeIdx = homeIdx
		teamsDup.VisitorIdx = visitorIdx

	}

	return teamsDup, nil
}
func findTeamByTeamIdx(teamsInLeague []storage.Team, teamIdxInLeague uint32) (*big.Int, error) {
	teamID := new(big.Int)
	i := sort.Search(len(teamsInLeague), func(i int) bool { return teamIdxInLeague <= teamsInLeague[i].TeamIdxInLeague })
	if i < len(teamsInLeague) && teamsInLeague[i].TeamIdxInLeague == teamIdxInLeague {
		teamID, _ = new(big.Int).SetString(teamsInLeague[i].TeamID, 10)
	} else {
		return teamID, errors.New("Not found teamId by team idx")
	}
	teamID, _ = new(big.Int).SetString(teamsInLeague[i].TeamID, 10)

	return teamID, nil
}

func (b *Calendar) Populate(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32, matchesStart [14]*big.Int, teamsInLeague []storage.Team) error {
	league, err := storage.LeagueByLeagueIdx(tx, leagueIdx)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Unexistent league")
	}
	if len(teamsInLeague) == 0 {
		log.Debugf("No teams provided fetching them...")
		teamsInLeague, err = storage.TeamsByTimezoneIdxCountryIdxLeagueIdxOrdered(tx, timezoneIdx, countryIdx, leagueIdx)
		if err != nil {
			return err
		}
	}
	var matchesToSetTeams []storage.Match
	for matchDay := uint8(0); matchDay < contracts.MatchDays; matchDay++ {
		for match := uint8(0); match < contracts.MatchesPerDay; match++ {
			teams, err := b.getTeamsInLeagueMatch(matchDay, match)
			if err != nil {
				return err
			}
			homeTeamID, err := findTeamByTeamIdx(teamsInLeague, uint32(teams.HomeIdx))
			if err != nil {
				return err
			}
			visitorTeamID, err := findTeamByTeamIdx(teamsInLeague, uint32(teams.VisitorIdx))
			if err != nil {
				return err
			}
			matchObj := storage.Match{
				TimezoneIdx:   timezoneIdx,
				CountryIdx:    countryIdx,
				LeagueIdx:     leagueIdx,
				MatchDayIdx:   matchDay,
				MatchIdx:      match,
				HomeTeamID:    homeTeamID,
				VisitorTeamID: visitorTeamID,
				HomeGoals:     0,
				VisitorGoals:  0,
				State:         "begin",
				StateExtra:    "",
				StartEpoch:    matchesStart[matchDay].Int64(),
			}
			matchesToSetTeams = append(matchesToSetTeams, matchObj)
		}
	}
	err = storage.MatchesSetTeamsBulkInsertUpdate(matchesToSetTeams, tx)
	if err != nil {
		return err
	}
	return nil
}

func (b *Calendar) Reset(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) error {
	league, err := storage.LeagueByLeagueIdx(tx, leagueIdx)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Unexistent league")
	}

	var matchesToReset []storage.Match
	for matchDay := uint8(0); matchDay < contracts.MatchDays; matchDay++ {
		for match := uint8(0); match < contracts.MatchesPerDay; match++ {
			matchObj := storage.Match{
				TimezoneIdx:   timezoneIdx,
				CountryIdx:    countryIdx,
				LeagueIdx:     leagueIdx,
				MatchDayIdx:   matchDay,
				MatchIdx:      match,
				HomeTeamID:    nil,
				VisitorTeamID: nil,
				HomeGoals:     0,
				VisitorGoals:  0,
				State:         "begin",
				StateExtra:    "",
				StartEpoch:    0,
			}
			matchesToReset = append(matchesToReset, matchObj)
		}
	}

	err = storage.MatchesResetBulkInsertUpdate(matchesToReset, tx)
	if err != nil {
		return err
	}
	return nil
}

func (b *Calendar) ResetByTzCountry(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagues []storage.League) error {
	var matchesToReset []storage.Match

	for _, league := range leagues {
		for matchDay := uint8(0); matchDay < contracts.MatchDays; matchDay++ {
			for match := uint8(0); match < contracts.MatchesPerDay; match++ {
				matchObj := storage.Match{
					TimezoneIdx:   timezoneIdx,
					CountryIdx:    countryIdx,
					LeagueIdx:     league.LeagueIdx,
					MatchDayIdx:   matchDay,
					MatchIdx:      match,
					HomeTeamID:    nil,
					VisitorTeamID: nil,
					HomeGoals:     0,
					VisitorGoals:  0,
					State:         "begin",
					StateExtra:    "",
					StartEpoch:    0,
				}
				matchesToReset = append(matchesToReset, matchObj)
			}
		}
	}

	err := storage.MatchesResetBulkInsertUpdate(matchesToReset, tx)
	if err != nil {
		return err
	}
	return nil
}

func (b *Calendar) PopulateByTzCountry(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, matchesStart [14]*big.Int, teamsInTzCountry []storage.Team, leagues []storage.League) error {

	var matchesToSetTeams []storage.Match
	for _, league := range leagues {

		for matchDay := uint8(0); matchDay < contracts.MatchDays; matchDay++ {
			for match := uint8(0); match < contracts.MatchesPerDay; match++ {
				teams, err := b.getTeamsInLeagueMatch(matchDay, match)
				if err != nil {
					return err
				}
				homeTeamID, err := findTeamByTeamIdx(teamsInTzCountry, uint32(teams.HomeIdx))
				if err != nil {
					return err
				}
				visitorTeamID, err := findTeamByTeamIdx(teamsInTzCountry, uint32(teams.VisitorIdx))
				if err != nil {
					return err
				}

				matchObj := storage.Match{
					TimezoneIdx:   timezoneIdx,
					CountryIdx:    countryIdx,
					LeagueIdx:     league.LeagueIdx,
					MatchDayIdx:   matchDay,
					MatchIdx:      match,
					HomeTeamID:    homeTeamID,
					VisitorTeamID: visitorTeamID,
					HomeGoals:     0,
					VisitorGoals:  0,
					State:         "begin",
					StateExtra:    "",
					StartEpoch:    matchesStart[matchDay].Int64(),
				}
				matchesToSetTeams = append(matchesToSetTeams, matchObj)
			}
		}

	}
	err := storage.MatchesSetTeamsBulkInsertUpdate(matchesToSetTeams, tx)
	if err != nil {
		return err
	}
	return nil
}
