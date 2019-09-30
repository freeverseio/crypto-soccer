package process

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/contracts/updates"
	"github.com/freeverseio/crypto-soccer/go-synchronizer/storage"

	log "github.com/sirupsen/logrus"
)

type LeagueProcessor struct {
	engine  *engine.Engine
	leagues *leagues.Leagues
	storage *storage.Storage
}

func NewLeagueProcessor(engine *engine.Engine, leagues *leagues.Leagues, storage *storage.Storage) *LeagueProcessor {
	return &LeagueProcessor{engine, leagues, storage}
}

func (b *LeagueProcessor) Process(event updates.UpdatesActionsSubmission) error {
	day := event.Day
	turnInDay := event.TurnInDay
	timezoneIdx := event.TimeZone
	if timezoneIdx < 1 || timezoneIdx > 24 {
		return errors.New("Wront timezone " + string(timezoneIdx))
	}
	if (turnInDay > 1) ||
		(turnInDay == 1 && day != 1) ||
		(turnInDay == 0 && (day < 2 || day > 14)) {
		log.Warnf("[LeagueProcessor] Skipping timezone %v, day %v, turnInDay %v", timezoneIdx, day, turnInDay)
		return nil
	}
	log.Infof("[LeagueProcessor] Processing timezone %v, day %v, turnInDay %v", timezoneIdx, day, turnInDay)

	countryCount, err := b.storage.CountryInTimezoneCount(timezoneIdx)
	if err != nil {
		return err
	}
	for countryIdx := uint32(0); countryIdx < countryCount; countryIdx++ {
		leagueCount, err := b.storage.LeagueInCountryCount(timezoneIdx, countryIdx)
		if err != nil {
			return err
		}
		for leagueIdx := uint32(0); leagueIdx < leagueCount; leagueIdx++ {
			matches, err := b.storage.GetMatchesInDay(timezoneIdx, countryIdx, leagueIdx, day-1)
			if err != nil {
				return err
			}
			for matchIdx := 0; matchIdx < len(matches); matchIdx++ {
				match := matches[matchIdx]
				matchSeed := big.NewInt(4) // TODO ??? what's this
				states, err := b.getMatchTeamsState(match.HomeTeamID, match.VisitorTeamID)
				if err != nil {
					return nil
				}
				var tactics [2]*big.Int
				tactics[0] = big.NewInt(0)
				tactics[1] = big.NewInt(0)
				is2ndHalf := false
				isHomeStadium := false
				result, err := b.engine.PlayMatch(
					&bind.CallOpts{},
					matchSeed,
					states,
					tactics,
					is2ndHalf,
					isHomeStadium,
				)
				if err != nil {
					log.Fatal(err)
					return err
				}
				err = b.storage.MatchSetResult(timezoneIdx, countryIdx, leagueIdx, uint32(day-1), uint32(matchIdx), result[0], result[1])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (b *LeagueProcessor) getMatchTeamsState(homeTeamID *big.Int, visitorTeamID *big.Int) ([2][25]*big.Int, error) {
	var states [2][25]*big.Int
	homeTeamState, err := b.getTeamState(homeTeamID)
	if err != nil {
		return states, err
	}
	visitorTeamState, err := b.getTeamState(visitorTeamID)
	if err != nil {
		return states, err
	}
	states[0] = homeTeamState
	states[1] = visitorTeamState
	return states, nil
}

func (b *LeagueProcessor) getTeamState(teamID *big.Int) ([25]*big.Int, error) {
	var state [25]*big.Int
	for i := 0; i < 25; i++ {
		state[i] = big.NewInt(435253)
	}
	return state, nil
}
