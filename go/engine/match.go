package engine

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"

	log "github.com/sirupsen/logrus"
)

type Match struct {
	contracts         *contracts.Contracts
	Seed              [32]byte
	StartTime         *big.Int
	HomeTeam          *Team
	VisitorTeam       *Team
	HomeGoals         uint8
	VisitorGoals      uint8
	HomeMatchLog      *big.Int
	VisitorMatchLog   *big.Int
	NOOUTOFGAMEPLAYER uint8
	REDCARD           uint8
	SOFTINJURY        uint8
	HARDINJURY        uint8
}

func (b Match) DumpState() string {
	var state string
	state += fmt.Sprintf("Seed: %v\n", b.Seed)
	state += fmt.Sprintf("StartTime: %v\n", b.StartTime)
	state += fmt.Sprintf("HomeTeam: %v\n", b.HomeTeam.DumpState())
	state += fmt.Sprintf("VisitorTeam: %v\n", b.VisitorTeam.DumpState())
	state += fmt.Sprintf("HomeGoals: %v\n", b.HomeGoals)
	state += fmt.Sprintf("VisitorGoals: %v\n", b.VisitorGoals)
	state += fmt.Sprintf("HomeMatchLog: %v\n", b.HomeMatchLog)
	state += fmt.Sprintf("VisitorMatchLog: %v\n", b.VisitorMatchLog)
	return state
}

func NewMatch(contracts *contracts.Contracts) (*Match, error) {
	var err error
	var mp Match
	mp.contracts = contracts
	mp.StartTime = big.NewInt(0)
	if mp.HomeTeam, err = NewTeam(contracts); err != nil {
		return nil, err
	}
	if mp.VisitorTeam, err = NewTeam(contracts); err != nil {
		return nil, err
	}
	mp.HomeMatchLog = big.NewInt(0)
	mp.VisitorMatchLog = big.NewInt(0)
	if mp.NOOUTOFGAMEPLAYER, err = contracts.Engineprecomp.NOOUTOFGAMEPLAYER(&bind.CallOpts{}); err != nil {
		return nil, err
	}
	if mp.REDCARD, err = contracts.Engineprecomp.REDCARD(&bind.CallOpts{}); err != nil {
		return nil, err
	}
	if mp.SOFTINJURY, err = contracts.Engineprecomp.SOFTINJURY(&bind.CallOpts{}); err != nil {
		return nil, err
	}
	if mp.HARDINJURY, err = contracts.Engineprecomp.HARDINJURY(&bind.CallOpts{}); err != nil {
		return nil, err
	}
	return &mp, nil
}

func (b *Match) Play1stHalf() error {
	isHomeStadium := true
	isPlayoff := false
	is2ndHalf := false
	matchSeed, err := b.generateMatchSeed()
	if err != nil {
		return err
	}
	matchLogs, err := b.contracts.Engine.PlayHalfMatch(
		&bind.CallOpts{},
		matchSeed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{b.HomeTeam.tactic, b.VisitorTeam.tactic},
		[2]*big.Int{b.HomeMatchLog, b.VisitorMatchLog},
		[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
	)
	if err != nil {
		return err
	}
	b.HomeMatchLog = matchLogs[0]
	b.VisitorMatchLog = matchLogs[1]
	goalsHome, goalsVisitor, err := b.getGoals(matchLogs)
	if err != nil {
		return err
	}
	b.HomeGoals += goalsHome
	b.VisitorGoals += goalsVisitor

	if err = b.updateTeamState(is2ndHalf, b.HomeTeam, b.HomeMatchLog); err != nil {
		return err
	}
	if err = b.updateTeamState(is2ndHalf, b.VisitorTeam, b.VisitorMatchLog); err != nil {
		return err
	}
	return nil
}

func (b *Match) Play2ndHalf() error {
	isHomeStadium := true
	isPlayoff := false
	is2ndHalf := true
	matchSeed, err := b.generateMatchSeed()
	if err != nil {
		return err
	}
	logs, err := b.contracts.Evolution.Play2ndHalfAndEvolve(
		&bind.CallOpts{},
		matchSeed,
		b.StartTime,
		b.Skills(),
		[2]*big.Int{b.HomeTeam.tactic, b.VisitorTeam.tactic},
		[2]*big.Int{b.HomeMatchLog, b.VisitorMatchLog},
		[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
	)
	if err != nil {
		return err
	}
	goalsHome, goalsVisitor, err := b.getGoals(logs)
	if err != nil {
		return err
	}
	b.HomeGoals += goalsHome
	b.VisitorGoals += goalsVisitor
	b.HomeMatchLog = logs[0]
	b.VisitorMatchLog = logs[1]
	if err = b.updateTeamState(is2ndHalf, b.HomeTeam, logs[0]); err != nil {
		return err
	}
	if err = b.updateTeamState(is2ndHalf, b.VisitorTeam, logs[1]); err != nil {
		return err
	}
	err = b.updateTeamSkills(b.HomeTeam, logs[0])
	if err != nil {
		return err
	}
	// err = b.updateTeamSkills(b.VisitorTeam, b.StartTime, logs[1])
	// if err != nil {
	// 	return err
	// }
	// err = b.updateTeamLeaderBoard()
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (b *Match) Skills() [2][25]*big.Int {
	return [2][25]*big.Int{b.HomeTeam.Skills(), b.VisitorTeam.Skills()}
}

func (b *Match) generateMatchSeed() (*big.Int, error) {
	matchSeed, err := b.contracts.Engine.GenerateMatchSeed(&bind.CallOpts{}, b.Seed, b.HomeTeam.TeamID, b.VisitorTeam.TeamID)
	if err != nil {
		return nil, err
	}
	z := new(big.Int)
	z.SetBytes(matchSeed[:])
	return z, nil
}

func (b *Match) ProcessMatchEvents(is2ndHalf bool) ([]storage.MatchEvent, error) {
	log.Warning("Match.ProcessMatchEvents TODO not implemented")
	return []storage.MatchEvent{}, nil
	// isHomeStadium := true
	// isPlayoff := false
	// states := b.Skills()
	// matchSeed, err := b.GenerateMatchSeed()
	// if err != nil {
	// 	return nil, err
	// }
	// seedAndStartTimeAndEvents, err := b.contracts.Matchevents.PlayHalfMatch(
	// 	&bind.CallOpts{},
	// 	matchSeed,
	// 	b.StartTime,
	// 	b.Skills(),
	// 	[2]*big.Int{b.HomeTeam.tactic, b.VisitorTeam.tactic},
	// 	[2]*big.Int{b.homeMatchLog, b.visitorMatchLog},
	// 	[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// events := seedAndStartTimeAndEvents[:]
	// log0, err := b.contracts.Utilsmatchlog.FullDecodeMatchLog(&bind.CallOpts{}, seedAndStartTimeAndEvents[0], is2ndHalf)
	// if err != nil {
	// 	return nil, err
	// }
	// log1, err := b.contracts.Utilsmatchlog.FullDecodeMatchLog(&bind.CallOpts{}, seedAndStartTimeAndEvents[1], is2ndHalf)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Debugf("Full decoded match log 0: %v", log0)
	// log.Debugf("Full decoded match log 1: %v", log1)
	// decodedTactics0, err := b.contracts.Assets.DecodeTactics(&bind.CallOpts{}, b.HomeTeam.tactic)
	// if err != nil {
	// 	return nil, err
	// }
	// decodedTactics1, err := b.contracts.Assets.DecodeTactics(&bind.CallOpts{}, b.VisitorTeam.tactic)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Debugf("Decoded tactics 0: %v", decodedTactics0)
	// log.Debugf("Decoded tactics 1: %v", decodedTactics1)
	// computedEvents, err := matchevents.GenerateMatchEvents(
	// 	matchSeed,
	// 	log0,
	// 	log1,
	// 	events,
	// 	decodedTactics0.Lineup,
	// 	decodedTactics1.Lineup,
	// 	decodedTactics0.Substitutions,
	// 	decodedTactics1.Substitutions,
	// 	decodedTactics0.SubsRounds,
	// 	decodedTactics1.SubsRounds,
	// 	is2ndHalf,
	// )
	// if err != nil {
	// 	return nil, err
	// }
	// var me []storage.MatchEvent
	// for _, computedEvent := range computedEvents {
	// 	var teamID string
	// 	if computedEvent.Team == 0 {
	// 		teamID = b.HomeTeam.TeamID.String()
	// 	} else if computedEvent.Team == 1 {
	// 		teamID = b.VisitorTeam.TeamID.String()
	// 	} else {
	// 		return nil, fmt.Errorf("Wrong match event team %v", computedEvent.Team)
	// 	}
	// 	event := storage.MatchEvent{}
	// 	event.TimezoneIdx = int(b.Match.TimezoneIdx)
	// 	event.CountryIdx = int(b.Match.CountryIdx)
	// 	event.LeagueIdx = int(b.Match.LeagueIdx)
	// 	event.MatchDayIdx = int(b.Match.MatchDayIdx)
	// 	event.MatchIdx = int(b.Match.MatchIdx)
	// 	event.TeamID = teamID
	// 	event.Minute = int(computedEvent.Minute)
	// 	event.Type, err = storage.MarchEventTypeByMatchEvent(computedEvent.Type)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	event.ManageToShoot = computedEvent.ManagesToShoot
	// 	event.IsGoal = computedEvent.IsGoal
	// 	primaryPlayerState := states[computedEvent.Team][computedEvent.PrimaryPlayer]
	// 	primaryPlayerID, err := b.contracts.Leagues.GetPlayerIdFromSkills(&bind.CallOpts{}, primaryPlayerState)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	event.PrimaryPlayerID = primaryPlayerID.String()
	// 	if computedEvent.SecondaryPlayer >= 0 && computedEvent.SecondaryPlayer < 25 {
	// 		secondaryPlayerState := states[computedEvent.Team][computedEvent.SecondaryPlayer]
	// 		secondaryPlayerID, err := b.contracts.Leagues.GetPlayerIdFromSkills(&bind.CallOpts{}, secondaryPlayerState)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		event.SecondaryPlayerID.String = secondaryPlayerID.String()
	// 		event.SecondaryPlayerID.Valid = true
	// 	}
	// 	me = append(me, event)
	// }
	// return me, nil
}

func (b *Match) getGoals(logs [2]*big.Int) (homeGoals uint8, VisitorGoals uint8, err error) {
	homeGoals, err = b.contracts.Evolution.GetNGoals(
		&bind.CallOpts{},
		logs[0],
	)
	if err != nil {
		return homeGoals, VisitorGoals, err
	}
	VisitorGoals, err = b.contracts.Evolution.GetNGoals(
		&bind.CallOpts{},
		logs[1],
	)
	return homeGoals, VisitorGoals, err
}

func (b *Match) updateTeamState(
	is2ndHalf bool,
	team *Team,
	matchLog *big.Int,
) error {
	decodedTactic, err := b.contracts.Leagues.DecodeTactics(&bind.CallOpts{}, team.tactic)
	if err != nil {
		return err
	}
	outOfGamePlayer, err := b.contracts.Engineprecomp.GetOutOfGamePlayer(&bind.CallOpts{}, matchLog, is2ndHalf)
	if err != nil {
		return err
	}
	outOfGameType, err := b.contracts.Engineprecomp.GetOutOfGameType(&bind.CallOpts{}, matchLog, is2ndHalf)
	if err != nil {
		return err
	}
	for shirt, player := range team.Players {
		wasAligned, err := b.contracts.Engine.WasPlayerAlignedEndOfLastHalf(
			&bind.CallOpts{},
			uint8(shirt),
			team.tactic,
			matchLog,
		)
		if err != nil {
			return err
		}
		player.skills, err = b.contracts.Evolution.SetAlignedEndOfLastHalf(
			&bind.CallOpts{},
			player.skills,
			wasAligned,
		)
		if err != nil {
			return err
		}
		var redCardMatchesLeft uint8
		var injuryMatchesLeft uint8
		if outOfGamePlayer.Int64() != int64(b.NOOUTOFGAMEPLAYER) {
			if outOfGamePlayer.Int64() < 0 || int(outOfGamePlayer.Int64()) >= len(decodedTactic.Lineup) {
				return fmt.Errorf("out of game player unknown %v, tactics %v, matchlog %v", outOfGamePlayer.Int64(), team.tactic, matchLog)
			}
			if uint8(shirt) == decodedTactic.Lineup[outOfGamePlayer.Int64()] {
				switch outOfGameType.Int64() {
				case int64(b.REDCARD):
					redCardMatchesLeft = 2
				case int64(b.SOFTINJURY):
					injuryMatchesLeft = 3
				case int64(b.HARDINJURY):
					injuryMatchesLeft = 7
				default:
					return fmt.Errorf("out of game type unknown %v", outOfGameType)
				}
			}
		}
		if is2ndHalf {
			if redCardMatchesLeft > 0 {
				redCardMatchesLeft--
			}
			if injuryMatchesLeft > 0 {
				injuryMatchesLeft--
			}
		}
		// log.Infof("encoded skills %v, redCard %v, injuries %v", player.EncodedSkills, player.RedCardMatchesLeft, player.InjuryMatchesLeft)
		if player.skills, err = b.contracts.Evolution.SetRedCardLastGame(&bind.CallOpts{}, player.skills, redCardMatchesLeft != 0); err != nil {
			return err
		}
		if player.skills, err = b.contracts.Evolution.SetInjuryWeeksLeft(&bind.CallOpts{}, player.skills, injuryMatchesLeft); err != nil {
			return err
		}

	}
	return nil
}

func (b *Match) updateTeamLeaderBoard() error {
	log.Warning("TODO Match::updateTeamLeaderBoard uninplemented")
	return nil

	// b.HomeTeam.GoalsForward += uint32(*b.Match.HomeGoals)
	// b.HomeTeam.GoalsAgainst += uint32(*b.Match.VisitorGoals)
	// b.VisitorTeam.GoalsForward += uint32(*b.Match.VisitorGoals)
	// b.VisitorTeam.GoalsAgainst += uint32(*b.Match.HomeGoals)

	// deltaGoals := int(*b.Match.HomeGoals) - int(*b.Match.VisitorGoals)
	// if deltaGoals > 0 {
	// 	b.HomeTeam.W++
	// 	b.VisitorTeam.L++
	// 	b.HomeTeam.Points += 3
	// } else if deltaGoals < 0 {
	// 	b.HomeTeam.L++
	// 	b.VisitorTeam.W++
	// 	b.VisitorTeam.Points += 3
	// } else {
	// 	b.HomeTeam.D++
	// 	b.VisitorTeam.D++
	// 	b.HomeTeam.Points++
	// 	b.VisitorTeam.Points++
	// }

	// return nil
}

func (b Match) getTrainingPoints() (homePoints, visitorPoints *big.Int, err error) {
	homePoints, err = b.contracts.Evolution.GetTrainingPoints(&bind.CallOpts{}, b.HomeMatchLog)
	if err != nil {
		return nil, nil, err
	}
	visitorPoints, err = b.contracts.Evolution.GetTrainingPoints(&bind.CallOpts{}, b.VisitorMatchLog)
	if err != nil {
		return nil, nil, err
	}
	return homePoints, visitorPoints, nil
}

func (b *Match) updateTeamSkills(
	team *Team,
	logs *big.Int,
) error {
	trainingPoints, err := b.contracts.Evolution.GetTrainingPoints(&bind.CallOpts{}, logs)
	if err != nil {
		return err
	}
	team.TrainingPoints = trainingPoints.Uint64()

	userAssignment, _ := new(big.Int).SetString("1022963800726800053580157736076735226208686447456863237", 10)
	newSkills, err := b.contracts.Evolution.GetTeamEvolvedSkills(
		&bind.CallOpts{},
		team.Skills(),
		userAssignment,
		b.StartTime,
	)
	if err != nil {
		return err
	}

	for i, player := range team.Players {
		player.skills = newSkills[i]
	}
	// TODO the followign code is old but it contain the code to evolve the name of a change in the generation of a player. Add it in the future.
	// for s, state := range newStates {
	// 	if state.String() == "0" {
	// 		continue
	// 	}

	// 	playerID, err := b.contracts.Leagues.GetPlayerIdFromSkills(&bind.CallOpts{}, state)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	player, err := storage.PlayerByPlayerId(tx, playerID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if player == nil {
	// 		return fmt.Errorf("Unexistent playerId %v", playerID)
	// 	}
	// 	oldGen, err := b.contracts.Assets.GetGeneration(&bind.CallOpts{}, states[s])
	// 	if err != nil {
	// 		return err
	// 	}
	// 	newGen, err := b.contracts.Assets.GetGeneration(&bind.CallOpts{}, state)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if newGen.Cmp(oldGen) != 0 {
	// 		timezone, countryIdx, _, err := b.contracts.Assets.DecodeTZCountryAndVal(&bind.CallOpts{}, player.PlayerId)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		newName, err := b.namesdb.GeneratePlayerFullName(player.PlayerId, uint8(newGen.Uint64()), timezone, countryIdx.Uint64())
	// 		if err != nil {
	// 			return err
	// 		}
	// 		player.Name = newName
	// 	}
	// 	defence, speed, pass, shoot, endurance, _, _, err := utils.DecodeSkills(b.contracts.Assets, state)
	// 	player.Defence = defence.Uint64()
	// 	player.Speed = speed.Uint64()
	// 	player.Pass = pass.Uint64()
	// 	player.Shoot = shoot.Uint64()
	// 	player.Defence = endurance.Uint64()
	// 	player.EncodedSkills = state
	// 	err = player.Update(tx)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
