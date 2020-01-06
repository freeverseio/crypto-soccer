package matchprocessor

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/matchevents"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/utils"

	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"

	log "github.com/sirupsen/logrus"
)

type Match struct {
	contracts          *contracts.Contracts
	Seed               [32]byte
	StartTime          *big.Int
	Match              *storage.Match
	HomeTeam           *storage.Team
	VisitorTeam        *storage.Team
	HomeTeamPlayers    []*storage.Player
	VisitorTeamPlayers []*storage.Player
	HomeTeamTactic     *big.Int
	VisitorTeamTactic  *big.Int
	NOOUTOFGAMEPLAYER  uint8
	REDCARD            uint8
	SOFTINJURY         uint8
	HARDINJURY         uint8
}

func NewMatch(
	contracts *contracts.Contracts,
	seed [32]byte,
	startTime *big.Int,
	match *storage.Match,
	homeTeam *storage.Team,
	visitorTeam *storage.Team,
	homeTeamPlayers []*storage.Player,
	visitorTeamPlayers []*storage.Player,
	homeTeamTactic *big.Int,
	visitorTeamTactic *big.Int,
) *Match {
	var mp Match
	mp.contracts = contracts
	mp.Seed = seed
	mp.StartTime = startTime
	mp.HomeTeamTactic = homeTeamTactic
	mp.VisitorTeamTactic = visitorTeamTactic
	mp.HomeTeam = homeTeam
	mp.VisitorTeam = visitorTeam
	mp.HomeTeamPlayers = homeTeamPlayers
	mp.VisitorTeamPlayers = visitorTeamPlayers
	mp.Match = match
	return &mp
}

func (b *Match) Process(is2ndHalf bool) ([]storage.MatchEvent, error) {
	var err error
	tactics := [2]*big.Int{b.HomeTeamTactic, b.VisitorTeamTactic} // TODO .Tactics() method

	var logs [2]*big.Int
	if is2ndHalf {
		logs, err = b.process2ndHalf()
	} else {
		logs, err = b.process1stHalf()
	}
	if err != nil {
		return nil, err
	}
	matchEvents, err := b.ProcessMatchEvents(is2ndHalf)
	if err != nil {
		return nil, err
	}
	goalsHome, goalsVisitor, err := b.GetGoals(logs)
	if err != nil {
		return nil, err
	}
	*b.Match.HomeGoals += goalsHome
	*b.Match.VisitorGoals += goalsVisitor
	b.Match.HomeMatchLog = logs[0]
	b.Match.VisitorMatchLog = logs[1]
	if err = b.UpdatePlayedByHalf(b.HomeTeamPlayers, is2ndHalf, b.HomeTeam.TeamID, tactics[0], logs[0]); err != nil {
		return nil, err
	}
	if err = b.UpdatePlayedByHalf(b.VisitorTeamPlayers, is2ndHalf, b.VisitorTeam.TeamID, tactics[1], logs[1]); err != nil {
		return nil, err
	}
	if is2ndHalf {
		err = b.UpdateTeamSkills(b.HomeTeam, b.HomeTeamPlayers, b.StartTime, logs[0])
		if err != nil {
			return nil, err
		}
		err = b.UpdateTeamSkills(b.VisitorTeam, b.VisitorTeamPlayers, b.StartTime, logs[1])
		if err != nil {
			return nil, err
		}
		err = b.updateTeamLeaderBoard()
		if err != nil {
			return nil, err
		}
	}
	return matchEvents, nil
}

func (b *Match) GenerateMatchSeed() (*big.Int, error) {
	matchSeed, err := b.contracts.Engine.GenerateMatchSeed(&bind.CallOpts{}, b.Seed, b.HomeTeam.TeamID, b.VisitorTeam.TeamID)
	if err != nil {
		return nil, err
	}
	z := new(big.Int)
	z.SetBytes(matchSeed[:])
	return z, nil
}

func (b *Match) ProcessMatchEvents(is2ndHalf bool) ([]storage.MatchEvent, error) {
	isHomeStadium := true
	isPlayoff := false
	matchSeed, err := b.GenerateMatchSeed()
	if err != nil {
		return nil, err
	}
	states, err := b.GetMatchTeamsState()
	if err != nil {
		return nil, err
	}
	seedAndStartTimeAndEvents, err := b.contracts.Matchevents.PlayHalfMatch(
		&bind.CallOpts{},
		matchSeed,
		b.StartTime,
		states,
		[2]*big.Int{b.HomeTeamTactic, b.VisitorTeamTactic},
		[2]*big.Int{b.Match.HomeMatchLog, b.Match.VisitorMatchLog},
		[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
	)
	if err != nil {
		return nil, err
	}

	events := seedAndStartTimeAndEvents[:]
	log0, err := b.contracts.Utilsmatchlog.FullDecodeMatchLog(&bind.CallOpts{}, seedAndStartTimeAndEvents[0], is2ndHalf)
	if err != nil {
		return nil, err
	}
	log1, err := b.contracts.Utilsmatchlog.FullDecodeMatchLog(&bind.CallOpts{}, seedAndStartTimeAndEvents[1], is2ndHalf)
	if err != nil {
		return nil, err
	}
	log.Debugf("Full decoded match log 0: %v", log0)
	log.Debugf("Full decoded match log 1: %v", log1)
	decodedTactics0, err := b.contracts.Assets.DecodeTactics(&bind.CallOpts{}, b.HomeTeamTactic)
	if err != nil {
		return nil, err
	}
	decodedTactics1, err := b.contracts.Assets.DecodeTactics(&bind.CallOpts{}, b.VisitorTeamTactic)
	if err != nil {
		return nil, err
	}
	log.Debugf("Decoded tactics 0: %v", decodedTactics0)
	log.Debugf("Decoded tactics 1: %v", decodedTactics1)
	computedEvents, err := matchevents.GenerateMatchEvents(
		matchSeed,
		log0,
		log1,
		events,
		decodedTactics0.Lineup,
		decodedTactics1.Lineup,
		decodedTactics0.Substitutions,
		decodedTactics1.Substitutions,
		decodedTactics0.SubsRounds,
		decodedTactics1.SubsRounds,
		is2ndHalf,
	)
	if err != nil {
		return nil, err
	}
	var me []storage.MatchEvent
	for _, computedEvent := range computedEvents {
		var teamID string
		if computedEvent.Team == 0 {
			teamID = b.HomeTeam.TeamID.String()
		} else if computedEvent.Team == 1 {
			teamID = b.VisitorTeam.TeamID.String()
		} else {
			return nil, fmt.Errorf("Wrong match event team %v", computedEvent.Team)
		}
		event := storage.MatchEvent{}
		event.TimezoneIdx = int(b.Match.TimezoneIdx)
		event.CountryIdx = int(b.Match.CountryIdx)
		event.LeagueIdx = int(b.Match.LeagueIdx)
		event.MatchDayIdx = int(b.Match.MatchDayIdx)
		event.MatchIdx = int(b.Match.MatchIdx)
		event.TeamID = teamID
		event.Minute = int(computedEvent.Minute)
		event.Type, err = storage.MarchEventTypeByMatchEvent(computedEvent.Type)
		if err != nil {
			return nil, err
		}
		event.ManageToShoot = computedEvent.ManagesToShoot
		event.IsGoal = computedEvent.IsGoal
		primaryPlayerState := states[computedEvent.Team][computedEvent.PrimaryPlayer]
		primaryPlayerID, err := b.contracts.Leagues.GetPlayerIdFromSkills(&bind.CallOpts{}, primaryPlayerState)
		if err != nil {
			return nil, err
		}
		event.PrimaryPlayerID = primaryPlayerID.String()
		if computedEvent.SecondaryPlayer >= 0 && computedEvent.SecondaryPlayer < 25 {
			secondaryPlayerState := states[computedEvent.Team][computedEvent.SecondaryPlayer]
			secondaryPlayerID, err := b.contracts.Leagues.GetPlayerIdFromSkills(&bind.CallOpts{}, secondaryPlayerState)
			if err != nil {
				return nil, err
			}
			event.SecondaryPlayerID.String = secondaryPlayerID.String()
			event.SecondaryPlayerID.Valid = true
		}
		me = append(me, event)
	}
	return me, nil
}

func (b *Match) process1stHalf() ([2]*big.Int, error) {
	var logs [2]*big.Int
	isHomeStadium := true
	isPlayoff := false
	is2ndHalf := false
	matchSeed, err := b.GenerateMatchSeed()
	if err != nil {
		return logs, err
	}
	states, err := b.GetMatchTeamsState()
	if err != nil {
		return logs, err
	}
	return b.contracts.Engine.PlayHalfMatch(
		&bind.CallOpts{},
		matchSeed,
		b.StartTime,
		states,
		[2]*big.Int{b.HomeTeamTactic, b.VisitorTeamTactic},
		[2]*big.Int{b.Match.HomeMatchLog, b.Match.VisitorMatchLog},
		[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
	)
}

func (b *Match) process2ndHalf() ([2]*big.Int, error) {
	var logs [2]*big.Int
	isHomeStadium := true
	isPlayoff := false
	is2ndHalf := true
	matchSeed, err := b.GenerateMatchSeed()
	if err != nil {
		return logs, err
	}
	states, err := b.GetMatchTeamsState()
	if err != nil {
		return logs, err
	}
	logs, err = b.contracts.Evolution.Play2ndHalfAndEvolve(
		&bind.CallOpts{},
		matchSeed,
		b.StartTime,
		states,
		[2]*big.Int{b.HomeTeamTactic, b.VisitorTeamTactic},
		[2]*big.Int{b.Match.HomeMatchLog, b.Match.VisitorMatchLog},
		[3]bool{is2ndHalf, isHomeStadium, isPlayoff},
	)
	return logs, err
}

func (b *Match) GetTeamState(players []*storage.Player) ([25]*big.Int, error) {
	var state [25]*big.Int
	for i := 0; i < 25; i++ {
		state[i] = big.NewInt(0)
	}
	for i := 0; i < len(players); i++ {
		player := players[i]
		playerSkills := player.EncodedSkills
		shirtNumber := player.ShirtNumber
		state[shirtNumber] = playerSkills
	}
	return state, nil
}

func (b *Match) GetMatchTeamsState() ([2][25]*big.Int, error) {
	var states [2][25]*big.Int
	homeTeamState, err := b.GetTeamState(b.HomeTeamPlayers)
	if err != nil {
		return states, err
	}
	visitorTeamState, err := b.GetTeamState(b.VisitorTeamPlayers)
	if err != nil {
		return states, err
	}
	states[0] = homeTeamState
	states[1] = visitorTeamState
	return states, nil
}
func (b *Match) GetGoals(logs [2]*big.Int) (homeGoals uint8, visitorGoals uint8, err error) {
	homeGoals, err = b.contracts.Evolution.GetNGoals(
		&bind.CallOpts{},
		logs[0],
	)
	if err != nil {
		return homeGoals, visitorGoals, err
	}
	visitorGoals, err = b.contracts.Evolution.GetNGoals(
		&bind.CallOpts{},
		logs[1],
	)
	return homeGoals, visitorGoals, err
}

func (b *Match) UpdatePlayedByHalf(
	players []*storage.Player,
	is2ndHalf bool,
	teamID *big.Int,
	tactic *big.Int,
	matchLog *big.Int,
) error {
	decodedTactic, err := b.contracts.Leagues.DecodeTactics(&bind.CallOpts{}, tactic)
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
	for _, player := range players {
		wasAligned, err := b.contracts.Engine.WasPlayerAlignedEndOfLastHalf(
			&bind.CallOpts{},
			player.ShirtNumber,
			tactic,
			matchLog,
		)
		if err != nil {
			return err
		}
		player.EncodedSkills, err = b.contracts.Evolution.SetAlignedEndOfLastHalf(
			&bind.CallOpts{},
			player.EncodedSkills,
			wasAligned,
		)
		if err != nil {
			return err
		}
		if outOfGamePlayer.Int64() != int64(b.NOOUTOFGAMEPLAYER) {
			if outOfGamePlayer.Int64() < 0 || int(outOfGamePlayer.Int64()) >= len(decodedTactic.Lineup) {
				return fmt.Errorf("out of game player unknown %v, tactics %v, matchlog %v", outOfGamePlayer.Int64(), tactic, matchLog)
			}
			if player.ShirtNumber == decodedTactic.Lineup[outOfGamePlayer.Int64()] {
				switch outOfGameType.Int64() {
				case int64(b.REDCARD):
					player.RedCardMatchesLeft = 2
				case int64(b.SOFTINJURY):
					player.InjuryMatchesLeft = 3
				case int64(b.HARDINJURY):
					player.InjuryMatchesLeft = 7
				default:
					return fmt.Errorf("out of game type unknown %v", outOfGameType)
				}
			}
		}
		if is2ndHalf {
			if player.RedCardMatchesLeft > 0 {
				player.RedCardMatchesLeft--
			}
			if player.InjuryMatchesLeft > 0 {
				player.InjuryMatchesLeft--
			}
		}
		// log.Infof("encoded skills %v, redCard %v, injuries %v", player.EncodedSkills, player.RedCardMatchesLeft, player.InjuryMatchesLeft)
		if player.EncodedSkills, err = b.contracts.Evolution.SetRedCardLastGame(&bind.CallOpts{}, player.EncodedSkills, player.RedCardMatchesLeft != 0); err != nil {
			return err
		}
		if player.EncodedSkills, err = b.contracts.Evolution.SetInjuryWeeksLeft(&bind.CallOpts{}, player.EncodedSkills, player.InjuryMatchesLeft); err != nil {
			return err
		}

	}
	return nil
}

func (b *Match) updateTeamLeaderBoard() error {
	b.HomeTeam.GoalsForward += uint32(*b.Match.HomeGoals)
	b.HomeTeam.GoalsAgainst += uint32(*b.Match.VisitorGoals)
	b.VisitorTeam.GoalsForward += uint32(*b.Match.VisitorGoals)
	b.VisitorTeam.GoalsAgainst += uint32(*b.Match.HomeGoals)

	deltaGoals := int(*b.Match.HomeGoals) - int(*b.Match.VisitorGoals)
	if deltaGoals > 0 {
		b.HomeTeam.W++
		b.VisitorTeam.L++
		b.HomeTeam.Points += 3
	} else if deltaGoals < 0 {
		b.HomeTeam.L++
		b.VisitorTeam.W++
		b.VisitorTeam.Points += 3
	} else {
		b.HomeTeam.D++
		b.VisitorTeam.D++
		b.HomeTeam.Points++
		b.VisitorTeam.Points++
	}

	return nil
}

func (b *Match) UpdateTeamSkills(
	team *storage.Team,
	players []*storage.Player,
	matchStartTime *big.Int,
	logs *big.Int,
) error {
	trainingPoints, err := b.contracts.Evolution.GetTrainingPoints(&bind.CallOpts{}, logs)
	if err != nil {
		return err
	}
	team.TrainingPoints = uint32(trainingPoints.Uint64())

	states, err := b.GetTeamState(players)
	if err != nil {
		return err
	}

	userAssignment, _ := new(big.Int).SetString("1022963800726800053580157736076735226208686447456863237", 10)
	newStates, err := b.contracts.Evolution.GetTeamEvolvedSkills(
		&bind.CallOpts{},
		states,
		userAssignment,
		matchStartTime,
	)
	if err != nil {
		return err
	}

	for i := range players {
		shirtNumber := players[i].ShirtNumber
		newState := newStates[shirtNumber]
		defence, speed, pass, shoot, endurance, _, _, err := utils.DecodeSkills(b.contracts.Assets, newState)
		if err != nil {
			return err
		}
		players[i].Defence = defence.Uint64()
		players[i].Speed = speed.Uint64()
		players[i].Pass = pass.Uint64()
		players[i].Shoot = shoot.Uint64()
		players[i].Defence = endurance.Uint64()
		players[i].EncodedSkills = newState
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

func GetEncodedTacticAtVerse(contracts *contracts.Contracts, teamID *big.Int, verse uint64) (*big.Int, error) {
	tactic := relay.DefaultTactic(teamID.String())
	// if tactic, err := relay.TacticByTeamIDAndVerse(teamID.String(), verse); err != nil {
	// 	return nil, err
	// } else
	if encodedTactic, err := contracts.Engine.EncodeTactics(
		&bind.CallOpts{},
		[3]uint8{11, 11, 11}, // TODO
		[3]uint8{2, 3, 4},    // TODO
		[14]uint8{
			uint8(tactic.Shirt0),
			uint8(tactic.Shirt1),
			uint8(tactic.Shirt2),
			uint8(tactic.Shirt3),
			uint8(tactic.Shirt4),
			uint8(tactic.Shirt5),
			uint8(tactic.Shirt6),
			uint8(tactic.Shirt7),
			uint8(tactic.Shirt8),
			uint8(tactic.Shirt9),
			uint8(tactic.Shirt10),
			uint8(tactic.Shirt11),
			uint8(tactic.Shirt12),
			uint8(tactic.Shirt13),
		},
		[10]bool{
			false,
			false,
			false,
			false,
			false,
			false,
			false,
			false,
			false,
			false,
		},
		uint8(tactic.TacticID),
	); err != nil {
		return nil, err
	} else {
		return encodedTactic, nil
	}
}
