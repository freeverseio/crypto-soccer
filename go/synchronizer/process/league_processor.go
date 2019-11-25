package process

import (
	"errors"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go/contracts/engineprecomp"
	"github.com/freeverseio/crypto-soccer/go/contracts/evolution"
	"github.com/freeverseio/crypto-soccer/go/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go/contracts/updates"
	"github.com/freeverseio/crypto-soccer/go/names"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"

	log "github.com/sirupsen/logrus"
)

type LeagueProcessor struct {
	engine            *engine.Engine
	enginePreComp     *engineprecomp.Engineprecomp
	leagues           *leagues.Leagues
	evolution         *evolution.Evolution
	universedb        *storage.Storage
	relaydb           *relay.Storage
	assets            *assets.Assets
	calendarProcessor *Calendar
	matchProcessor    *MatchProcessor
}

func NewLeagueProcessor(
	engine *engine.Engine,
	enginePreComp *engineprecomp.Engineprecomp,
	assets *assets.Assets,
	leagues *leagues.Leagues,
	evolution *evolution.Evolution,
	universedb *storage.Storage,
	relaydb *relay.Storage,
	namesdb *names.Generator,
) (*LeagueProcessor, error) {
	calendarProcessor, err := NewCalendar(leagues, universedb)
	if err != nil {
		return nil, err
	}
	matchProcessor, err := NewMatchProcessor(
		universedb,
		relaydb,
		assets,
		leagues,
		evolution,
		engine,
		enginePreComp,
		namesdb,
	)
	if err != nil {
		return nil, err
	}
	return &LeagueProcessor{
		engine,
		enginePreComp,
		leagues,
		evolution,
		universedb,
		relaydb,
		assets,
		calendarProcessor,
		matchProcessor,
	}, nil
}

func (b *LeagueProcessor) Process(event updates.UpdatesActionsSubmission) error {
	day := event.Day
	turnInDay := event.TurnInDay
	timezoneIdx := event.TimeZone
	log.Debugf("[LeagueProcessor] Processing timezone %v, day %v, turnInDay %v", timezoneIdx, day, turnInDay)

	if timezoneIdx > 24 {
		return errors.New("[LaegueProcessor] ... wront timezone")
	}

	// switch turnInDay {
	// case 0: // first half league match
	// case 1:
	if turnInDay < 2 {
		countryCount, err := b.universedb.CountryInTimezoneCount(timezoneIdx)
		if err != nil {
			return err
		}
		for countryIdx := uint32(0); countryIdx < countryCount; countryIdx++ {
			// if a new league is starting shuffle the teams
			if (day == 0) && (turnInDay == 0) {
				err = b.UpdatePrevPerfPointsAndShuffleTeamsInCountry(timezoneIdx, countryIdx)
				if err != nil {
					return err
				}
			}
			leagueCount, err := b.universedb.LeagueInCountryCount(timezoneIdx, countryIdx)
			if err != nil {
				return err
			}
			for leagueIdx := uint32(0); leagueIdx < leagueCount; leagueIdx++ {
				if day == 0 {
					err = b.resetLeague(timezoneIdx, countryIdx, leagueIdx)
					if err != nil {
						return err
					}
				}
				matches, err := b.universedb.GetMatchesInDay(timezoneIdx, countryIdx, leagueIdx, day)
				if err != nil {
					return err
				}
				for _, match := range matches {
					is2ndHalf := turnInDay == 1
					err = b.matchProcessor.Process(
						match,
						event.Seed,
						event.SubmissionTime,
						is2ndHalf,
					)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	// default:
	// 	log.Warnf("[LeagueProcessor] ... skipping")
	// } // switch
	return nil
}

func (b *LeagueProcessor) UpdatePrevPerfPointsAndShuffleTeamsInCountry(timezoneIdx uint8, countryIdx uint32) error {
	log.Infof("[LeagueProcessor] Shuffling timezone %v, country %v", timezoneIdx, countryIdx)
	var orgMap []storage.Team
	leagueCount, err := b.universedb.LeagueInCountryCount(timezoneIdx, countryIdx)
	if err != nil {
		return err
	}
	for leagueIdx := uint32(0); leagueIdx < leagueCount; leagueIdx++ {
		teams, err := b.universedb.GetTeamsInLeague(timezoneIdx, countryIdx, leagueIdx)
		if err != nil {
			return err
		}
		// ordening by points
		sort.Slice(teams[:], func(i, j int) bool {
			return teams[i].State.Points > teams[j].State.Points
		})
		for position, team := range teams {
			teamState, err := b.matchProcessor.GetTeamState(team.TeamID)
			if err != nil {
				return err
			}
			if !storage.IsBotTeam(team) {
				team.State.RankingPoints, team.State.PrevPerfPoints, err = b.leagues.ComputeTeamRankingPoints(
					&bind.CallOpts{},
					teamState,
					uint8(position),
					team.State.PrevPerfPoints,
					team.TeamID,
				)
				if err != nil {
					return err
				}
			}
			// log.Infof("New ranking team %v points %v ranking %v", team.TeamID, team.State.Points, newPrevPerfPoints)
			orgMap = append(orgMap, team)
		}
	}
	// ordening all the teams by ranking points
	sort.Slice(orgMap[:], func(i, j int) bool {
		return orgMap[i].State.RankingPoints.Cmp(orgMap[j].State.RankingPoints) == 1
	})
	// create the new leagues
	for i, team := range orgMap {
		team.State.LeagueIdx = uint32(i / 8)
		team.State.TeamIdxInLeague = uint32(i % 8)
		err = b.universedb.TeamUpdate(team.TeamID, team.State)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *LeagueProcessor) resetLeague(timezoneIdx uint8, countryIdx uint32, leagueIdx uint32) error {
	teams, err := b.universedb.GetTeamsInLeague(timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		return err
	}
	for i := 0; i < len(teams); i++ {
		team := teams[i]
		team.State.D = 0
		team.State.W = 0
		team.State.L = 0
		team.State.GoalsAgainst = 0
		team.State.GoalsForward = 0
		team.State.Points = 0
		err = b.universedb.TeamUpdate(team.TeamID, team.State)
		if err != nil {
			return err
		}
	}
	err = b.calendarProcessor.Reset(timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		return err
	}
	err = b.calendarProcessor.Populate(timezoneIdx, countryIdx, leagueIdx)
	if err != nil {
		return err
	}
	return nil
}
