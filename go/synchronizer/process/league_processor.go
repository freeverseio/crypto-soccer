package process

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go/contracts/engineprecomp"
	"github.com/freeverseio/crypto-soccer/go/contracts/evolution"
	"github.com/freeverseio/crypto-soccer/go/contracts/leagues"
	"github.com/freeverseio/crypto-soccer/go/contracts/updates"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"

	log "github.com/sirupsen/logrus"
)

type LeagueProcessor struct {
	engine            *engine.Engine
	enginePreComp     *engineprecomp.Engineprecomp
	leagues           *leagues.Leagues
	evolution         *evolution.Evolution
	storage           *storage.Storage
	calendarProcessor *Calendar
	playerHackSkills  *big.Int
}

func NewLeagueProcessor(
	engine *engine.Engine,
	enginePreComp *engineprecomp.Engineprecomp,
	leagues *leagues.Leagues,
	evolution *evolution.Evolution,
	storage *storage.Storage,
) (*LeagueProcessor, error) {
	calendarProcessor, err := NewCalendar(leagues, storage)
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
		enginePreComp,
		leagues,
		evolution,
		storage,
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

	// switch turnInDay {
	// case 0: // first half league match
	// case 1:
	if turnInDay < 2 {
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
				if day == 0 {
					err = b.resetLeague(timezoneIdx, countryIdx, leagueIdx)
					if err != nil {
						return err
					}
				}
				matches, err := b.storage.GetMatchesInDay(timezoneIdx, countryIdx, leagueIdx, day)
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
					is2ndHalf := turnInDay == 1
					isHomeStadium := true
					isPlayoff := false
					var matchLog [2]*big.Int
					var matchBools [3]bool
					matchBools[0] = is2ndHalf
					matchBools[1] = isHomeStadium
					matchBools[2] = isPlayoff
					var logs [2]*big.Int
					if is2ndHalf {
						matchLog[0], matchLog[1], err = b.storage.GetMatchLogs(timezoneIdx, countryIdx, leagueIdx, day, uint8(matchIdx))
						if err != nil {
							return nil
						}
						logs, err = b.evolution.Play2ndHalfAndEvolve(
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
					} else { // first half
						matchLog[0] = big.NewInt(0)
						matchLog[1] = big.NewInt(0)
						logs, err = b.engine.PlayHalfMatch(
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
					err = b.storage.MatchSetResult(timezoneIdx, countryIdx, leagueIdx, day, uint8(matchIdx), goalsHome, goalsVisitor, logs[0], logs[1])
					if err != nil {
						return err
					}
					err = b.updateTeamStatistics(match.HomeTeamID, match.VisitorTeamID, goalsHome, goalsVisitor)
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
				}
			}
		}
	}
	// default:
	// 	log.Warnf("[LeagueProcessor] ... skipping")
	// } // switch
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
	players, err := b.storage.GetPlayersOfTeam(teamID)
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
					player.State.EncodedSkills, err = b.evolution.SetRedCardLastGame(&bind.CallOpts{}, player.State.EncodedSkills, true)
					if err != nil {
						return err
					}
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
			player.State.EncodedSkills, err = b.evolution.SetRedCardLastGame(&bind.CallOpts{}, player.State.EncodedSkills, false)
			if err != nil {
				return err
			}
		}
		b.storage.PlayerUpdate(player.PlayerId, player.State)
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
	teams, err := b.storage.GetTeamsInLeague(timezoneIdx, countryIdx, leagueIdx)
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
		err = b.storage.TeamUpdate(team.TeamID, team.State)
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
	homeTeam, err := b.storage.GetTeam(homeTeamID)
	if err != nil {
		return err
	}
	visitorTeam, err := b.storage.GetTeam(visitorTeamID)
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

	err = b.storage.TeamUpdate(homeTeamID, homeTeam.State)
	if err != nil {
		return err
	}
	err = b.storage.TeamUpdate(visitorTeamID, visitorTeam.State)
	if err != nil {
		return err
	}
	return nil
}

func (b *LeagueProcessor) GetMatchTactics(homeTeamID *big.Int, visitorTeamID *big.Int) ([2]*big.Int, error) {
	var tactics [2]*big.Int
	var substitutions [3]uint8 = [3]uint8{11, 11, 11} // TODO: get the constant from the contracts => 11 means no subs
	var subsRounds [3]uint8 = [3]uint8{2, 3, 4}
	var lineup [14]uint8 = [14]uint8{0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	var extraAttack [10]bool
	extraAttack[3] = true
	extraAttack[6] = true
	var tacticsId uint8 = 1
	tactic, err := b.engine.EncodeTactics(
		&bind.CallOpts{},
		substitutions,
		subsRounds,
		lineup,
		extraAttack,
		tacticsId,
	)
	if err != nil {
		return tactics, err
	}
	tactics[0] = tactic
	tactics[1] = tactic
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
	freePlayerID, err := b.engine.FREEPLAYERID(&bind.CallOpts{})
	if err != nil {
		return state, err
	}
	for i := 0; i < 25; i++ {
		state[i] = freePlayerID
	}
	players, err := b.storage.GetPlayersOfTeam(teamID)
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
