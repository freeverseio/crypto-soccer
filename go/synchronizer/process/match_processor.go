package process

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/matchevents"
	"github.com/freeverseio/crypto-soccer/go/names"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/utils"
	log "github.com/sirupsen/logrus"
)

type MatchProcessor struct {
	contracts         *contracts.Contracts
	namesdb           *names.Generator
	NOOUTOFGAMEPLAYER uint8
	REDCARD           uint8
	SOFTINJURY        uint8
	HARDINJURY        uint8
}

func NewMatchProcessor(
	contracts *contracts.Contracts,
	namesdb *names.Generator,
) (*MatchProcessor, error) {
	processor := MatchProcessor{}
	var err error
	processor.NOOUTOFGAMEPLAYER, err = contracts.Engineprecomp.NOOUTOFGAMEPLAYER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	processor.REDCARD, err = contracts.Engineprecomp.REDCARD(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	processor.SOFTINJURY, err = contracts.Engineprecomp.SOFTINJURY(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	processor.HARDINJURY, err = contracts.Engineprecomp.HARDINJURY(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	processor.contracts = contracts
	processor.namesdb = namesdb

	return &processor, nil
}

func (b *MatchProcessor) GetGoals(logs [2]*big.Int) (homeGoals uint8, visitorGoals uint8, err error) {
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

func (b *MatchProcessor) ProcessMatchEvents(
	match storage.Match,
	states [2][25]*big.Int,
	tactics [2]*big.Int,
	matchSeed *big.Int,
	startTime *big.Int,
	is2ndHalf bool,
) ([]storage.MatchEvent, error) {
	isHomeStadium := true
	isPlayoff := false
	matchLog := [2]*big.Int{}
	if is2ndHalf { // TODO make match.HomeMatchLog and Visitor to be 0 if first half
		matchLog = [2]*big.Int{big.NewInt(0), big.NewInt(0)}
	} else {
		matchLog = [2]*big.Int{match.HomeMatchLog, match.VisitorMatchLog}
	}
	matchBools := [3]bool{is2ndHalf, isHomeStadium, isPlayoff}
	seedAndStartTimeAndEvents, err := b.contracts.Matchevents.PlayHalfMatch(
		&bind.CallOpts{},
		matchSeed,
		startTime,
		states,
		tactics,
		matchLog,
		matchBools,
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
	decodedTactics0, err := b.contracts.Assets.DecodeTactics(&bind.CallOpts{}, tactics[0])
	if err != nil {
		return nil, err
	}
	decodedTactics1, err := b.contracts.Assets.DecodeTactics(&bind.CallOpts{}, tactics[1])
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
			teamID = match.HomeTeamID.String()
		} else if computedEvent.Team == 1 {
			teamID = match.VisitorTeamID.String()
		} else {
			return nil, fmt.Errorf("Wrong match event team %v", computedEvent.Team)
		}
		event := storage.MatchEvent{}
		event.TimezoneIdx = int(match.TimezoneIdx)
		event.CountryIdx = int(match.CountryIdx)
		event.LeagueIdx = int(match.LeagueIdx)
		event.MatchDayIdx = int(match.MatchDayIdx)
		event.MatchIdx = int(match.MatchIdx)
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

func (b *MatchProcessor) InnerProcess(
	seed [32]byte,
	startTime *big.Int,
	is2ndHalf bool,
	match *storage.Match,
	homeTeam *storage.Team,
	visitorTeam *storage.Team,
	homeTeamPlayers []*storage.Player,
	visitorTeamPlayers []*storage.Player,
) ([]storage.MatchEvent, error) {
	var err error
	tactics, err := b.GetMatchTactics(match.HomeTeamID, match.VisitorTeamID)
	if err != nil {
		return nil, err
	}
	matchSeed, err := b.GenerateMatchSeed(seed, match.HomeTeamID, match.VisitorTeamID)
	if err != nil {
		return nil, err
	}
	states, err := b.GetMatchTeamsState(homeTeamPlayers, visitorTeamPlayers)
	if err != nil {
		return nil, err
	}
	var logs [2]*big.Int
	if is2ndHalf {
		logs, err = b.process2ndHalf(*match, states, tactics, matchSeed, startTime)
	} else {
		logs, err = b.process1stHalf(*match, states, tactics, matchSeed, startTime)
	}
	if err != nil {
		return nil, err
	}
	matchEvents, err := b.ProcessMatchEvents(*match, states, tactics, matchSeed, startTime, is2ndHalf)
	if err != nil {
		return nil, err
	}
	goalsHome, goalsVisitor, err := b.GetGoals(logs)
	if err != nil {
		return nil, err
	}
	match.HomeGoals = &goalsHome
	match.VisitorGoals = &goalsVisitor
	match.HomeMatchLog = new(big.Int).Set(logs[0])
	match.VisitorMatchLog = new(big.Int).Set(logs[1])
	if err = b.UpdatePlayedByHalf(homeTeamPlayers, is2ndHalf, match.HomeTeamID, tactics[0], logs[0]); err != nil {
		return nil, err
	}
	if err = b.UpdatePlayedByHalf(visitorTeamPlayers, is2ndHalf, match.VisitorTeamID, tactics[1], logs[1]); err != nil {
		return nil, err
	}
	if is2ndHalf {
		err = b.UpdateTeamSkills(homeTeam, homeTeamPlayers, startTime, logs[0])
		if err != nil {
			return nil, err
		}
		err = b.UpdateTeamSkills(visitorTeam, visitorTeamPlayers, startTime, logs[1])
		if err != nil {
			return nil, err
		}
		err = b.updateTeamLeaderBoard(homeTeam, visitorTeam, goalsHome, goalsVisitor)
		if err != nil {
			return nil, err
		}
	}
	return matchEvents, nil
}

func (b *MatchProcessor) Process(
	tx *sql.Tx,
	match storage.Match,
	seed [32]byte,
	startTime *big.Int,
	is2ndHalf bool,
) error {
	log.Debugf("MatchProcessor::Process Tz: %v, Country: %v, league: %v, matchDay: %v, match: %v, 2ndHalf: %v",
		match.TimezoneIdx,
		match.CountryIdx,
		match.LeagueIdx,
		match.MatchDayIdx,
		match.MatchIdx,
		is2ndHalf,
	)
	homeTeam, err := storage.TeamByTeamId(tx, match.HomeTeamID)
	if err != nil {
		return err
	}
	visitorTeam, err := storage.TeamByTeamId(tx, match.VisitorTeamID)
	if err != nil {
		return err
	}
	homeTeamPlayers, err := storage.PlayersByTeamId(tx, match.HomeTeamID)
	if err != nil {
		return err
	}
	visitorTeamPlayers, err := storage.PlayersByTeamId(tx, match.VisitorTeamID)
	if err != nil {
		return err
	}
	events, err := b.InnerProcess(
		seed,
		startTime,
		is2ndHalf,
		&match,
		&homeTeam,
		&visitorTeam,
		homeTeamPlayers,
		visitorTeamPlayers,
	)
	if err != nil {
		return err
	}
	err = homeTeam.Update(tx)
	if err != nil {
		return err
	}
	err = visitorTeam.Update(tx)
	if err != nil {
		return err
	}
	if err = match.Update(tx); err != nil {
		return err
	}
	for _, player := range homeTeamPlayers {
		if err = player.Update(tx); err != nil {
			return err
		}
	}
	for _, player := range visitorTeamPlayers {
		if err = player.Update(tx); err != nil {
			return err
		}
	}
	for _, event := range events {
		if err = event.Insert(tx); err != nil {
			return err
		}
	}
	return nil
}

func (b *MatchProcessor) GetTeamState(players []*storage.Player) ([25]*big.Int, error) {
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

func (b *MatchProcessor) GenerateMatchSeed(seed [32]byte, homeTeamID *big.Int, visitorTeamID *big.Int) (*big.Int, error) {
	matchSeed, err := b.contracts.Engine.GenerateMatchSeed(&bind.CallOpts{}, seed, homeTeamID, visitorTeamID)
	if err != nil {
		return nil, err
	}
	z := new(big.Int)
	z.SetBytes(matchSeed[:])
	return z, nil
}
func (b *MatchProcessor) GetMatchTeamsState(
	homeTeamPlayers []*storage.Player,
	visitorTeamPlayers []*storage.Player,
) ([2][25]*big.Int, error) {
	var states [2][25]*big.Int
	homeTeamState, err := b.GetTeamState(homeTeamPlayers)
	if err != nil {
		return states, err
	}
	visitorTeamState, err := b.GetTeamState(visitorTeamPlayers)
	if err != nil {
		return states, err
	}
	states[0] = homeTeamState
	states[1] = visitorTeamState
	return states, nil
}

func (b *MatchProcessor) process1stHalf(
	match storage.Match,
	states [2][25]*big.Int,
	tactics [2]*big.Int,
	matchSeed *big.Int,
	startTime *big.Int,
) (logs [2]*big.Int, err error) {
	isHomeStadium := true
	isPlayoff := false
	is2ndHalf := false
	matchLog := [2]*big.Int{big.NewInt(0), big.NewInt(0)}
	matchBools := [3]bool{is2ndHalf, isHomeStadium, isPlayoff}
	return b.contracts.Engine.PlayHalfMatch(
		&bind.CallOpts{},
		matchSeed,
		startTime,
		states,
		tactics,
		matchLog,
		matchBools,
	)
}

func (b *MatchProcessor) process2ndHalf(
	match storage.Match,
	states [2][25]*big.Int,
	tactics [2]*big.Int,
	matchSeed *big.Int,
	startTime *big.Int,
) (logs [2]*big.Int, err error) {
	isHomeStadium := true
	isPlayoff := false
	is2ndHalf := true
	matchLog := [2]*big.Int{match.HomeMatchLog, match.VisitorMatchLog}
	matchBools := [3]bool{is2ndHalf, isHomeStadium, isPlayoff}
	logs, err = b.contracts.Evolution.Play2ndHalfAndEvolve(
		&bind.CallOpts{},
		matchSeed,
		startTime,
		states,
		tactics,
		matchLog,
		matchBools,
	)
	return logs, err
}

func (b *MatchProcessor) GetMatchTactics(homeTeamID *big.Int, visitorTeamID *big.Int) ([2]*big.Int, error) {
	var tactics [2]*big.Int
	verse := uint64(0) // TODO: get verse from event
	tactic, err := b.getEncodedTacticAtVerse(homeTeamID, verse)
	if err != nil {
		return tactics, err
	}
	tactics[0] = tactic
	tactic, err = b.getEncodedTacticAtVerse(visitorTeamID, verse)
	if err != nil {
		return tactics, err
	}
	tactics[1] = tactic
	return tactics, nil
}

func (b *MatchProcessor) UpdatePlayedByHalf(
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
				log.Warningf("out of game player unknown %v, tactics %v, matchlog %v", outOfGamePlayer.Int64(), tactic, matchLog)
				continue
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

func (b *MatchProcessor) updateTeamLeaderBoard(homeTeam *storage.Team, visitorTeam *storage.Team, homeGoals uint8, visitorGoals uint8) error {
	homeTeam.GoalsForward += uint32(homeGoals)
	homeTeam.GoalsAgainst += uint32(visitorGoals)
	visitorTeam.GoalsForward += uint32(visitorGoals)
	visitorTeam.GoalsAgainst += uint32(homeGoals)

	deltaGoals := int(homeGoals) - int(visitorGoals)
	if deltaGoals > 0 {
		homeTeam.W++
		visitorTeam.L++
		homeTeam.Points += 3
	} else if deltaGoals < 0 {
		homeTeam.L++
		visitorTeam.W++
		visitorTeam.Points += 3
	} else {
		homeTeam.D++
		visitorTeam.D++
		homeTeam.Points++
		visitorTeam.Points++
	}

	return nil
}

func (b *MatchProcessor) getEncodedTacticAtVerse(teamID *big.Int, verse uint64) (*big.Int, error) {
	tactic := relay.DefaultTactic(teamID.String())
	// if tactic, err := relay.TacticByTeamIDAndVerse(teamID.String(), verse); err != nil {
	// 	return nil, err
	// } else
	if encodedTactic, err := b.contracts.Engine.EncodeTactics(
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

func (b *MatchProcessor) UpdateTeamSkills(
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
