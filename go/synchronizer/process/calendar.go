package process

import (
	"database/sql"
	"errors"
	"math/big"

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

func (b *Calendar) Populate(tx *sql.Tx, timezoneIdx uint8, countryIdx uint32, leagueIdx uint32, matchesStart [14]*big.Int) error {
	league, err := storage.LeagueByLeagueIdx(tx, leagueIdx)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("Unexistent league")
	}
	var matchesToSetTeams []storage.Match
	for matchDay := uint8(0); matchDay < contracts.MatchDays; matchDay++ {
		for match := uint8(0); match < contracts.MatchesPerDay; match++ {
			teams, err := b.contracts.Leagues.GetTeamsInLeagueMatch(&bind.CallOpts{}, matchDay, match)
			if err != nil {
				return nil
			}
			homeTeamID, err := storage.TeamIdByTimezoneIdxCountryIdxLeagueIdx(tx, timezoneIdx, countryIdx, leagueIdx, uint32(teams.HomeIdx))
			if err != nil {
				return err
			}
			visitorTeamID, err := storage.TeamIdByTimezoneIdxCountryIdxLeagueIdx(tx, timezoneIdx, countryIdx, leagueIdx, uint32(teams.VisitorIdx))
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
			// err = storage.MatchSetTeams(tx, timezoneIdx, countryIdx, leagueIdx, matchDay, match, homeTeamID, visitorTeamID, matchesStart[matchDay])
			// if err != nil {
			// 	return err
			// }
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
	log.Infof("Num matchesToReset %v", len(matchesToReset))

	err = storage.MatchesResetBulkInsertUpdate(matchesToReset, tx)
	if err != nil {
		return err
	}
	return nil
}
