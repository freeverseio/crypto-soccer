package process

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts/assets"
	"github.com/freeverseio/crypto-soccer/go/contracts/engine"
	"github.com/freeverseio/crypto-soccer/go/contracts/engineprecomp"
	"github.com/freeverseio/crypto-soccer/go/contracts/evolution"
	"github.com/freeverseio/crypto-soccer/go/contracts/leagues"
	relay "github.com/freeverseio/crypto-soccer/go/relay/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/utils"
)

type MatchProcessor struct {
	universedb        *storage.Storage
	relaydb           *relay.Storage
	assets            *assets.Assets
	leagues           *leagues.Leagues
	evolution         *evolution.Evolution
	engine            *engine.Engine
	enginePreComp     *engineprecomp.Engineprecomp
	NOOUTOFGAMEPLAYER uint8
	REDCARD           uint8
	SOFTINJURY        uint8
	HARDINJURY        uint8
}

func NewMatchProcessor(
	universedb *storage.Storage,
	relaydb *relay.Storage,
	assets *assets.Assets,
	leagues *leagues.Leagues,
	evolution *evolution.Evolution,
	engine *engine.Engine,
	enginePreComp *engineprecomp.Engineprecomp,
) (*MatchProcessor, error) {
	processor := MatchProcessor{}
	var err error
	processor.NOOUTOFGAMEPLAYER, err = enginePreComp.NOOUTOFGAMEPLAYER(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	processor.REDCARD, err = enginePreComp.REDCARD(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	processor.SOFTINJURY, err = enginePreComp.SOFTINJURY(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	processor.HARDINJURY, err = enginePreComp.HARDINJURY(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	processor.universedb = universedb
	processor.relaydb = relaydb
	processor.assets = assets
	processor.leagues = leagues
	processor.evolution = evolution
	processor.engine = engine
	processor.enginePreComp = enginePreComp
	return &processor, nil
}

func (b *MatchProcessor) GetGoals(logs [2]*big.Int) (homeGoals uint8, visitorGoals uint8, err error) {
	homeGoals, err = b.evolution.GetNGoals(
		&bind.CallOpts{},
		logs[0],
	)
	if err != nil {
		return homeGoals, visitorGoals, err
	}
	visitorGoals, err = b.evolution.GetNGoals(
		&bind.CallOpts{},
		logs[1],
	)
	return homeGoals, visitorGoals, err
}

func (b *MatchProcessor) Process(
	match storage.Match,
	seed [32]byte,
	startTime *big.Int,
	is2ndHalf bool,
) error {
	tactics, err := b.GetMatchTactics(match.HomeTeamID, match.VisitorTeamID)
	if err != nil {
		return err
	}
	states, err := b.GetMatchTeamsState(match.HomeTeamID, match.VisitorTeamID)
	if err != nil {
		return err
	}
	// play the match half
	var logs [2]*big.Int
	if is2ndHalf {
		logs, err = b.process2ndHalf(match, states, tactics, seed, startTime)
	} else {
		logs, err = b.process1stHalf(match, states, tactics, seed, startTime)
	}
	if err != nil {
		return err
	}
	goalsHome, goalsVisitor, err := b.GetGoals(logs)
	if err != nil {
		return err
	}
	err = b.universedb.MatchSetResult(
		match.TimezoneIdx,
		match.CountryIdx,
		match.LeagueIdx,
		match.MatchDayIdx,
		match.MatchIdx,
		goalsHome,
		goalsVisitor,
		logs[0],
		logs[1],
	)
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
		err = b.UpdateTeamSkills(states[0], startTime, logs[0])
		if err != nil {
			return err
		}
		err = b.UpdateTeamSkills(states[1], startTime, logs[1])
		if err != nil {
			return err
		}
		err = b.updateTeamLeaderBoard(match.HomeTeamID, match.VisitorTeamID, goalsHome, goalsVisitor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *MatchProcessor) GetTeamState(teamID *big.Int) ([25]*big.Int, error) {
	var state [25]*big.Int
	for i := 0; i < 25; i++ {
		state[i] = big.NewInt(0)
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

func (b *MatchProcessor) GenerateMatchSeed(seed [32]byte, homeTeamID *big.Int, visitorTeamID *big.Int) (*big.Int, error) {
	matchSeed, err := b.engine.GenerateMatchSeed(&bind.CallOpts{}, seed, homeTeamID, visitorTeamID)
	if err != nil {
		return nil, err
	}
	z := new(big.Int)
	z.SetBytes(matchSeed[:])
	return z, nil
}
func (b *MatchProcessor) GetMatchTeamsState(homeTeamID *big.Int, visitorTeamID *big.Int) ([2][25]*big.Int, error) {
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

func (b *MatchProcessor) process1stHalf(
	match storage.Match,
	states [2][25]*big.Int,
	tactics [2]*big.Int,
	seed [32]byte,
	startTime *big.Int,
) (logs [2]*big.Int, err error) {
	matchSeed, err := b.GenerateMatchSeed(seed, match.HomeTeamID, match.VisitorTeamID)
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

func (b *MatchProcessor) process2ndHalf(
	match storage.Match,
	states [2][25]*big.Int,
	tactics [2]*big.Int,
	seed [32]byte,
	startTime *big.Int,
) (logs [2]*big.Int, err error) {
	matchSeed, err := b.GenerateMatchSeed(seed, match.HomeTeamID, match.VisitorTeamID)
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

func (b *MatchProcessor) UpdatePlayedByHalf(is2ndHalf bool, teamID *big.Int, tactic *big.Int, matchLog *big.Int) error {
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
	for _, player := range players {
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
		if outOfGamePlayer.Int64() != int64(b.NOOUTOFGAMEPLAYER) {
			if player.State.ShirtNumber == decodedTactic.Lineup[outOfGamePlayer.Int64()] {
				switch outOfGameType.Int64() {
				case int64(b.REDCARD):
					player.State.RedCardMatchesLeft = 2
				case int64(b.SOFTINJURY):
					player.State.InjuryMatchesLeft = 3
				case int64(b.HARDINJURY):
					player.State.InjuryMatchesLeft = 7
				}
				if player.State.EncodedSkills, err = b.evolution.SetRedCardLastGame(&bind.CallOpts{}, player.State.EncodedSkills, player.State.RedCardMatchesLeft == 0); err != nil {
					return err
				}
				if player.State.EncodedSkills, err = b.evolution.SetInjuryWeeksLeft(&bind.CallOpts{}, player.State.EncodedSkills, player.State.InjuryMatchesLeft); err != nil {
					return err
				}
			}
		}
		if err = b.universedb.PlayerUpdate(player.PlayerId, player.State); err != nil {
			return nil
		}
	}
	return nil
}

func (b *MatchProcessor) updateTeamLeaderBoard(homeTeamID *big.Int, visitorTeamID *big.Int, homeGoals uint8, visitorGoals uint8) error {
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

func (b *MatchProcessor) getEncodedTacticAtVerse(teamID *big.Int, verse uint64) (*big.Int, error) {
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

func (b *MatchProcessor) UpdateTeamSkills(
	states [25]*big.Int,
	matchStartTime *big.Int,
	logs *big.Int,
) error {
	trainingPoints, err := b.evolution.GetTrainingPoints(&bind.CallOpts{}, logs)
	if err != nil {
		return err
	}
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
		if state.String() == "0" {
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
