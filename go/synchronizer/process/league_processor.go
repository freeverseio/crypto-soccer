package process

import (
	"errors"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go/contracts/engineprecomp"
	"github.com/freeverseio/crypto-soccer/go/contracts/evolution"
	"github.com/freeverseio/crypto-soccer/go/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go/contracts/updates"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/utils"

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
	FREEPLAYERID      *big.Int
}

func NewLeagueProcessor(
	engine *engine.Engine,
	enginePreComp *engineprecomp.Engineprecomp,
	assets *assets.Assets,
	leagues *leagues.Leagues,
	evolution *evolution.Evolution,
	universedb *storage.Storage,
	relaydb *relay.Storage,
) (*LeagueProcessor, error) {
	calendarProcessor, err := NewCalendar(leagues, universedb)
	if err != nil {
		return nil, err
	}

	FREEPLAYERID, err := engine.FREEPLAYERID(&bind.CallOpts{})
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
		FREEPLAYERID,
	}, nil
}

func (b *LeagueProcessor) process1stHalf(
	match storage.Match,
	tactics [2]*big.Int,
	seed [32]byte,
	startTime *big.Int,
) (logs [2]*big.Int, err error) {
	matchSeed, err := b.GenerateMatchSeed(seed, match.HomeTeamID, match.VisitorTeamID)
	if err != nil {
		return logs, err
	}
	states, err := b.GetMatchTeamsState(match.HomeTeamID, match.VisitorTeamID)
	if err != nil {
		return logs, err
	}

	isHomeStadium := true
	isPlayoff := false
	is2ndHalf := false
	matchLog := [2]*big.Int{big.NewInt(0), big.NewInt(0)}
	matchBools := [3]bool{is2ndHalf, isHomeStadium, isPlayoff}
	return b.engine.PlayHalfMatch(
		&bind.CallOpts{},
		matchSeed,
		startTime,
		states,
		tactics,
		matchLog,
		matchBools,
	)
}

func (b *LeagueProcessor) process2ndHalf(
	match storage.Match,
	tactics [2]*big.Int,
	seed [32]byte,
	startTime *big.Int,
) (logs [2]*big.Int, err error) {
	matchSeed, err := b.GenerateMatchSeed(seed, match.HomeTeamID, match.VisitorTeamID)
	if err != nil {
		return logs, err
	}
	states, err := b.GetMatchTeamsState(match.HomeTeamID, match.VisitorTeamID)
	if err != nil {
		return logs, err
	}
	isHomeStadium := true
	isPlayoff := false
	is2ndHalf := true
	matchLog := [2]*big.Int{match.HomeMatchLog, match.VisitorMatchLog}
	matchBools := [3]bool{is2ndHalf, isHomeStadium, isPlayoff}
	logs, err = b.evolution.Play2ndHalfAndEvolve(
		&bind.CallOpts{},
		matchSeed,
		startTime,
		states,
		tactics,
		matchLog,
		matchBools,
	)
	if err != nil {
		return logs, err
	}
	for i := 0; i < 2; i++ {
		trainingPointHomeTeam, err := b.evolution.GetTrainingPoints(&bind.CallOpts{}, logs[i])
		if err != nil {
			return logs, err
		}
		err = b.UpdateTeamSkills(states[i], trainingPointHomeTeam, startTime)
		if err != nil {
			return logs, err
		}
	}
	return logs, err
}

func (b *LeagueProcessor) processHalfMatch(
	match storage.Match,
	tactics [2]*big.Int,
	seed [32]byte,
	startTime *big.Int,
	is2ndHalf bool,
) (logs [2]*big.Int, err error) {
	if is2ndHalf {
		return b.process2ndHalf(match, tactics, seed, startTime)
	}
	return b.process1stHalf(match, tactics, seed, startTime)
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
			if (day == 0) && (turnInDay == 0) {
				err = b.shuffleTeamsInCountry(timezoneIdx, countryIdx)
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
				for matchIdx := 0; matchIdx < len(matches); matchIdx++ {
					match := matches[matchIdx]
					is2ndHalf := turnInDay == 1
					tactics, err := b.GetMatchTactics(match.HomeTeamID, match.VisitorTeamID)
					if err != nil {
						return err
					}
					logs, err := b.processHalfMatch(match, tactics, event.Seed, event.SubmissionTime, is2ndHalf)
					if err != nil {
						return err
					}
					goalsHome, err := b.evolution.GetNGoals(
						&bind.CallOpts{},
						logs[0],
					)
					if err != nil {
						return err
					}
					goalsVisitor, err := b.evolution.GetNGoals(
						&bind.CallOpts{},
						logs[1],
					)
					if err != nil {
						return err
					}
					err = b.universedb.MatchSetResult(timezoneIdx, countryIdx, leagueIdx, day, uint8(matchIdx), goalsHome, goalsVisitor, logs[0], logs[1])
					if err != nil {
						return err
					}
					err = b.UpdatePlayedByHalf(is2ndHalf, match.HomeTeamID, tactics[0], logs[0])
					if err != nil {
						return err
					}
					err = b.UpdatePlayedByHalf(is2ndHalf, match.VisitorTeamID, tactics[1], logs[1])
					if err != nil {
						return err
					}
					if is2ndHalf {
						err = b.updateTeamLeaderBoard(match.HomeTeamID, match.VisitorTeamID, goalsHome, goalsVisitor)
						if err != nil {
							return err
						}
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

func (b *LeagueProcessor) shuffleTeamsInCountry(timezoneIdx uint8, countryIdx uint32) error {
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
			teamState, err := b.GetTeamState(team.TeamID)
			if err != nil {
				return err
			}
			team.State.RankingPoints, err = b.leagues.ComputeTeamRankingPoints(
				&bind.CallOpts{},
				teamState,
				uint8(position),
				team.State.RankingPoints,
			)
			if err != nil {
				return err
			}
			// log.Infof("New ranking team %v points %v ranking %v", team.TeamID, team.State.Points, newRankingPoints)
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

func (b *LeagueProcessor) UpdateTeamSkills(states [25]*big.Int, trainingPoints *big.Int, matchStartTime *big.Int) error {
	userAssignment, _ := new(big.Int).SetString("1022963800726800053580157736076735226208686447456863237", 10)
	newStates, err := b.evolution.GetTeamEvolvedSkills(
		&bind.CallOpts{},
		states,
		trainingPoints,
		userAssignment,
		matchStartTime,
	)
	if err != nil {
		return err
	}

	for _, state := range newStates {
		if state.String() == b.FREEPLAYERID.String() {
			continue
		}

		playerID, err := b.leagues.GetPlayerIdFromSkills(&bind.CallOpts{}, state)
		if err != nil {
			return err
		}
		player, err := b.universedb.GetPlayer(playerID)
		if err != nil {
			return err
		}
		defence, speed, pass, shoot, endurance, _, _, err := utils.DecodeSkills(b.assets, state)
		player.State.Defence = defence.Uint64()
		player.State.Speed = speed.Uint64()
		player.State.Pass = pass.Uint64()
		player.State.Shoot = shoot.Uint64()
		player.State.Defence = endurance.Uint64()
		player.State.EncodedSkills = state
		err = b.universedb.PlayerUpdate(playerID, player.State)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *LeagueProcessor) UpdatePlayedByHalf(is2ndHalf bool, teamID *big.Int, tactic *big.Int, matchLog *big.Int) error {
	NO_OUT_OF_GAME_PLAYER, err := b.enginePreComp.NOOUTOFGAMEPLAYER(&bind.CallOpts{})
	if err != nil {
		return err
	}
	RED_CARD, err := b.enginePreComp.REDCARD(&bind.CallOpts{})
	if err != nil {
		return err
	}
	SOFTINJURY, err := b.enginePreComp.SOFTINJURY(&bind.CallOpts{})
	if err != nil {
		return err
	}
	HARDINJURY, err := b.enginePreComp.HARDINJURY(&bind.CallOpts{})
	if err != nil {
		return err
	}
	players, err := b.universedb.GetPlayersOfTeam(teamID)
	if err != nil {
		return err
	}
	decodedTactic, err := b.leagues.DecodeTactics(&bind.CallOpts{}, tactic)
	if err != nil {
		return err
	}
	outOfGamePlayer, err := b.enginePreComp.GetOutOfGamePlayer(&bind.CallOpts{}, matchLog, is2ndHalf)
	if err != nil {
		return err
	}
	outOfGameType, err := b.enginePreComp.GetOutOfGameType(&bind.CallOpts{}, matchLog, is2ndHalf)
	if err != nil {
		return err
	}
	for i := 0; i < len(players); i++ {
		player := players[i]
		wasAligned, err := b.engine.WasPlayerAlignedEndOfLastHalf(
			&bind.CallOpts{},
			player.State.ShirtNumber,
			tactic,
			matchLog,
		)
		if err != nil {
			return err
		}
		player.State.EncodedSkills, err = b.evolution.SetAlignedEndOfLastHalf(
			&bind.CallOpts{},
			player.State.EncodedSkills,
			wasAligned,
		)
		if err != nil {
			return err
		}
		if outOfGamePlayer.Int64() != int64(NO_OUT_OF_GAME_PLAYER) {
			if player.State.ShirtNumber == decodedTactic.Lineup[outOfGamePlayer.Int64()] {
				switch outOfGameType.Int64() {
				case int64(RED_CARD):
					player.State.RedCard = true
				case int64(SOFTINJURY):
					player.State.EncodedSkills, err = b.evolution.SetInjuryWeeksLeft(&bind.CallOpts{}, player.State.EncodedSkills, 1)
					if err != nil {
						return err
					}
				case int64(HARDINJURY):
					player.State.EncodedSkills, err = b.evolution.SetInjuryWeeksLeft(&bind.CallOpts{}, player.State.EncodedSkills, 2)
					if err != nil {
						return err
					}
				}
			}
		}
		if is2ndHalf {
			player.State.RedCard = false
		}
		player.State.EncodedSkills, err = b.evolution.SetRedCardLastGame(&bind.CallOpts{}, player.State.EncodedSkills, player.State.RedCard)
		if err != nil {
			return err
		}
		if err = b.universedb.PlayerUpdate(player.PlayerId, player.State); err != nil {
			return nil
		}
	}
	return nil
}

func (b *LeagueProcessor) GenerateMatchSeed(seed [32]byte, homeTeamID *big.Int, visitorTeamID *big.Int) (*big.Int, error) {
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

func (b *LeagueProcessor) updateTeamLeaderBoard(homeTeamID *big.Int, visitorTeamID *big.Int, homeGoals uint8, visitorGoals uint8) error {
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
	if tactic, err := b.relaydb.GetTacticOrDefault(teamID, verse); err != nil {
		return nil, err
	} else if encodedTactic, err := b.engine.EncodeTactics(
		&bind.CallOpts{},
		tactic.Substitutions,
		tactic.SubsRounds,
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
		state[i] = b.FREEPLAYERID
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
