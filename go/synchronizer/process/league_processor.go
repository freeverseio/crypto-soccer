package process

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go/contracts/evolution"
	"github.com/freeverseio/crypto-soccer/go/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go/contracts/updates"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"

	log "github.com/sirupsen/logrus"
)

type LeagueProcessor struct {
	engine            *engine.Engine
	leagues           *leagues.Leagues
	evolution         *evolution.Evolution
	universedb        *storage.Storage
	relaydb           *relay.Storage
	calendarProcessor *Calendar
	playerHackSkills  *big.Int
}

func NewLeagueProcessor(
	engine *engine.Engine,
	leagues *leagues.Leagues,
	evolution *evolution.Evolution,
	universedb *storage.Storage,
	relaydb *relay.Storage,
) (*LeagueProcessor, error) {
	calendarProcessor, err := NewCalendar(leagues, universedb)
	if err != nil {
		return nil, err
	}

	playerHackSkills, _ := new(big.Int).SetString("713624055286353394965726120199142814938406092850", 10)
	if err != nil {
		return nil, err
	}
	// playerHackSkills := big.NewInt(0)

	return &LeagueProcessor{
		engine,
		leagues,
		evolution,
		universedb,
		relaydb,
		calendarProcessor,
		playerHackSkills,
	}, nil
}

func (b *LeagueProcessor) Process(event updates.UpdatesActionsSubmission) error {
	day := event.Day
	turnInDay := event.TurnInDay
	timezoneIdx := event.TimeZone
	log.Infof("[LeagueProcessor] Processing timezone %v, day %v, turnInDay %v", timezoneIdx, day, turnInDay)

	if timezoneIdx > 24 {
		return errors.New("[LaegueProcessor] ... wront timezone")
	}
	isFirstHalfLeagueMatch := turnInDay == 0
	if isFirstHalfLeagueMatch == false {
		log.Warnf("[LeagueProcessor] ... skipping")
		return nil
	}

	countryCount, err := b.universedb.CountryInTimezoneCount(timezoneIdx)
	if err != nil {
		return err
	}
	for countryIdx := uint32(0); countryIdx < countryCount; countryIdx++ {
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
			for matchIdx := 0; matchIdx < len(matches); matchIdx++ {
				match := matches[matchIdx]
				matchSeed, err := b.GenerateMatchSeed(event.Seed, match.HomeTeamID, match.VisitorTeamID)
				if err != nil {
					return err
				}
				states, err := b.GetMatchTeamsState(match.HomeTeamID, match.VisitorTeamID)
				if err != nil {
					return err
				}
				tactics, err := b.GetMatchTactics(match.HomeTeamID, match.VisitorTeamID)
				if err != nil {
					return err
				}
				is2ndHalf := false
				isHomeStadium := true
				isPlayoff := false
				var matchLog [2]*big.Int
				matchLog[0] = big.NewInt(0)
				matchLog[1] = big.NewInt(0)
				var matchBools [3]bool
				matchBools[0] = is2ndHalf
				matchBools[1] = isHomeStadium
				matchBools[2] = isPlayoff
				result, err := b.engine.PlayHalfMatch(
					&bind.CallOpts{},
					matchSeed,
					event.SubmissionTime,
					states,
					tactics,
					matchLog,
					matchBools,
				)
				if err != nil {
					return err
				}
				goalsHome, err := b.evolution.GetNGoals(
					&bind.CallOpts{},
					result[0],
				)
				if err != nil {
					return err
				}
				goalsVisitor, err := b.evolution.GetNGoals(
					&bind.CallOpts{},
					result[1],
				)
				if err != nil {
					return err
				}
				err = b.universedb.MatchSetResult(timezoneIdx, countryIdx, leagueIdx, uint32(day), uint32(matchIdx), goalsHome, goalsVisitor)
				if err != nil {
					return err
				}
				err = b.updateTeamStatistics(match.HomeTeamID, match.VisitorTeamID, goalsHome, goalsVisitor)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (b *LeagueProcessor) GenerateMatchSeed(seed [32]byte, homeTeamID *big.Int, visitorTeamID *big.Int) (*big.Int, error) {
	// uint256Ty, _ := abi.NewType("uint256", nil)
	// bytes32Ty, _ := abi.NewType("bytes32", nil)

	// arguments := abi.Arguments{
	// 	{
	// 		Type: bytes32Ty,
	// 	},
	// 	{
	// 		Type: uint256Ty,
	// 	},
	// 	{
	// 		Type: uint256Ty,
	// 	},
	// }

	// bytes, _ := arguments.Pack(
	// 	seed,
	// 	homeTeamID,
	// 	visitorTeamID,
	// )
	// crypto.Keccak256(bytes)

	matchSeed, err := b.engine.GenerateMatchSeed(&bind.CallOpts{}, seed, homeTeamID, visitorTeamID)
	if err != nil {
		return nil, err
	}
	z := new(big.Int)
	z.SetBytes(matchSeed[:])
	return z, nil
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

func (b *LeagueProcessor) updateTeamStatistics(homeTeamID *big.Int, visitorTeamID *big.Int, homeGoals uint8, visitorGoals uint8) error {
	homeTeam, err := b.universedb.GetTeam(homeTeamID)
	if err != nil {
		return err
	}
	visitorTeam, err := b.universedb.GetTeam(visitorTeamID)
	if err != nil {
		return err
	}

	homeTeam.State.GoalsForward += uint32(homeGoals)
	homeTeam.State.GoalsAgainst += uint32(visitorGoals)
	visitorTeam.State.GoalsForward += uint32(visitorGoals)
	visitorTeam.State.GoalsAgainst += uint32(homeGoals)

	deltaGoals := int(homeGoals) - int(visitorGoals)
	if deltaGoals > 0 {
		homeTeam.State.W++
		visitorTeam.State.L++
		homeTeam.State.Points += 3
	} else if deltaGoals < 0 {
		homeTeam.State.L++
		visitorTeam.State.W++
		visitorTeam.State.Points += 3
	} else {
		homeTeam.State.D++
		visitorTeam.State.D++
		homeTeam.State.Points++
		visitorTeam.State.Points++
	}

	err = b.universedb.TeamUpdate(homeTeamID, homeTeam.State)
	if err != nil {
		return err
	}
	err = b.universedb.TeamUpdate(visitorTeamID, visitorTeam.State)
	if err != nil {
		return err
	}
	return nil
}

func (b *LeagueProcessor) getEncodedTacticAtVerse(teamID *big.Int, verse uint64) (*big.Int, error) {
	substitutions := [3]uint8{11, 11, 11} // 11 no substitutions // TODO (future)
	subsRounds := [3]uint8{2, 3, 4}       // TODO (future)
	if tactic, err := b.relaydb.GetTacticOrDefault(teamID, verse); err != nil {
		return nil, err
	} else if encodedTactic, err := b.engine.EncodeTactics(
		&bind.CallOpts{},
		substitutions,
		subsRounds,
		tactic.Shirts,
		tactic.ExtraAttack,
		tactic.TacticID,
	); err != nil {
		return nil, err
	} else {
		return encodedTactic, nil
	}
}

func (b *LeagueProcessor) GetMatchTactics(homeTeamID *big.Int, visitorTeamID *big.Int) ([2]*big.Int, error) {
	var tactics [2]*big.Int
	verse := uint64(0) // TODO: get verse from event
	if tactic, err := b.getEncodedTacticAtVerse(homeTeamID, verse); err != nil {
		return tactics, err
	} else {
		tactics[0] = tactic
	}
	if tactic, err := b.getEncodedTacticAtVerse(visitorTeamID, verse); err != nil {
		return tactics, err
	} else {
		tactics[1] = tactic
	}
	return tactics, nil
}

func (b *LeagueProcessor) GetMatchTeamsState(homeTeamID *big.Int, visitorTeamID *big.Int) ([2][25]*big.Int, error) {
	var states [2][25]*big.Int
	homeTeamState, err := b.GetTeamState(homeTeamID)
	if err != nil {
		return states, err
	}
	visitorTeamState, err := b.GetTeamState(visitorTeamID)
	if err != nil {
		return states, err
	}
	states[0] = homeTeamState
	states[1] = visitorTeamState
	return states, nil
}

func (b *LeagueProcessor) GetTeamState(teamID *big.Int) ([25]*big.Int, error) {
	var state [25]*big.Int
	for i := 0; i < 25; i++ {
		state[i] = b.playerHackSkills
	}
	players, err := b.universedb.GetPlayersOfTeam(teamID)
	if err != nil {
		return state, err
	}
	for i := 0; i < len(players); i++ {
		player := players[i]
		playerSkills := player.State.EncodedSkills
		shirtNumber := player.State.ShirtNumber
		state[shirtNumber] = playerSkills
	}
	return state, nil
}
