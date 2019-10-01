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

	day-- // cause we use 0 starting indexes

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
			matches, err := b.storage.GetMatchesInDay(timezoneIdx, countryIdx, leagueIdx, day)
			if err != nil {
				return err
			}
			for matchIdx := 0; matchIdx < len(matches); matchIdx++ {
				match := matches[matchIdx]
				matchSeed := big.NewInt(4) // TODO ??? what's this
				states, err := b.getMatchTeamsState(match.HomeTeamID, match.VisitorTeamID)
				if err != nil {
					log.Error(err)
					return err
				}
				tactics, err := b.getMatchTactics(match.HomeTeamID, match.VisitorTeamID)
				if err != nil {
					return err
				}
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
				err = b.setResult(timezoneIdx, countryIdx, leagueIdx, day, matchIdx, result[0], result[1])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (b *LeagueProcessor) setResult(
	timezoneIdx uint8,
	countryIdx uint32,
	leagueIdx uint32,
	day uint8,
	matchIdx int,
	homeGoals uint8,
	visitorGoals uint8,
) error {
	err := b.storage.MatchSetResult(timezoneIdx, countryIdx, leagueIdx, uint32(day), uint32(matchIdx), homeGoals, visitorGoals)
	if err != nil {
		return err
	}
	return nil
}

func (b *LeagueProcessor) getMatchTactics(homeTeamID *big.Int, visitorTeamID *big.Int) ([2]*big.Int, error) {
	var tactics [2]*big.Int
	tactics[0], _ = new(big.Int).SetString("1216069450684002467840", 10)
	tactics[1], _ = new(big.Int).SetString("1216069450684002467840", 10)
	return tactics, nil
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
		// playerSkills, err := b.leagues.GetPlayerSkillsAtBirth(&bind.CallOpts{}, big.NewInt(1)) // TODO remove 1 with something better
		// if err != nil {
		// 	return state, err
		// }
		state[i] = big.NewInt(0)
	}
	players, err := b.storage.GetPlayersOfTeam(teamID)
	if err != nil {
		return state, err
	}
	for i := 0; i < len(players); i++ {
		playerID := players[i].PlayerId
		playerSkills, err := b.leagues.GetPlayerSkillsAtBirth(&bind.CallOpts{}, playerID)
		if err != nil {
			return state, err
		}
		state[i] = playerSkills
	}
	return state, nil
}
