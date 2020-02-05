package engine

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/storage"
)

const NoOutOfGamePlayer = uint8(14)
const RedCard = uint8(3)
const SoftInjury = uint8(1)
const HardInjury = uint8(2)

type Team struct {
	storage.Team
	Players [25]*Player
	tactic  *big.Int // TODO add to storage.Team
}

func NewTeam() *Team {
	var team Team
	team.TeamID = big.NewInt(0)
	for i := range team.Players {
		team.Players[i] = NewPlayer()
	}
	team.tactic = DefaultTactic()
	return &team
}

func (b Team) ToStorage(contracts contracts.Contracts, tx *sql.Tx) error {
	for _, player := range b.Players {
		stoPlayer, err := player.ToStorage(contracts)
		if err != nil {
			return err
		}
		if err := stoPlayer.Update(tx); err != nil {
			return err
		}
	}
	return b.Update(tx)
}

func (b Team) DumpState() string {
	var state string
	state += fmt.Sprintf("TeamId: %v\n", b.TeamID)
	state += fmt.Sprintf("Points: %v\n", b.Points)
	state += fmt.Sprintf("W: %v\n", b.W)
	state += fmt.Sprintf("D: %v\n", b.D)
	state += fmt.Sprintf("L: %v\n", b.L)
	state += fmt.Sprintf("GoalsForward: %v\n", b.GoalsForward)
	state += fmt.Sprintf("GoalsAgainst: %v\n", b.GoalsAgainst)
	for i, player := range b.Players {
		state += fmt.Sprintf("Players[%d]: %v\n", i, player.DumpState())
	}
	state += fmt.Sprintf("tactic: %v\n", b.tactic)
	state += fmt.Sprintf("TrainingPoints: %v", b.TrainingPoints)
	return state
}

func (b Team) Skills() [25]*big.Int {
	var skills [25]*big.Int
	for i := range skills {
		skills[i] = b.Players[i].Skills()
	}
	return skills
}

func DefaultTactic() *big.Int {
	tactic, _ := new(big.Int).SetString("340596594427581673436941882753025", 10)
	return tactic
}

func (b *Team) SetSkills(contracts contracts.Contracts, skills [25]*big.Int) {
	for i := range skills {
		b.Players[i].SetSkills(skills[i])
	}
}

// func (b *Team) Evolve(
// 	contracts contracts.Contracts,
// 	matchLog *big.Int,
// 	startTime *big.Int,
// 	is2ndHalf bool,
// ) error {
// 	if err := b.updateTeamState(contracts, is2ndHalf, matchLog); err != nil {
// 		return err
// 	}
// 	if is2ndHalf {
// 		if err := b.updateTeamSkills(contracts, matchLog, startTime); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (b *Team) updateTeamState(
// 	contracts contracts.Contracts,
// 	is2ndHalf bool,
// 	matchLog *big.Int,
// ) error {
// 	decodedTactic, err := contracts.Leagues.DecodeTactics(&bind.CallOpts{}, b.tactic)
// 	if err != nil {
// 		return err
// 	}
// 	outOfGamePlayer, err := contracts.Engineprecomp.GetOutOfGamePlayer(&bind.CallOpts{}, matchLog, is2ndHalf)
// 	if err != nil {
// 		return err
// 	}
// 	outOfGameType, err := contracts.Engineprecomp.GetOutOfGameType(&bind.CallOpts{}, matchLog, is2ndHalf)
// 	if err != nil {
// 		return err
// 	}
// 	for shirt, player := range b.Players {
// 		wasAligned, err := contracts.Engine.WasPlayerAlignedEndOfLastHalf(
// 			&bind.CallOpts{},
// 			uint8(shirt),
// 			b.tactic,
// 			matchLog,
// 		)
// 		if err != nil {
// 			return err
// 		}
// 		if err = player.SetAligned(contracts, wasAligned); err != nil {
// 			return err
// 		}
// 		var redCardMatchesLeft uint8
// 		var injuryMatchesLeft uint8
// 		if outOfGamePlayer.Int64() != int64(NoOutOfGamePlayer) {
// 			if uint8(shirt) == decodedTactic.Lineup[outOfGamePlayer.Int64()] {
// 				switch outOfGameType.Int64() {
// 				case int64(RedCard):
// 					redCardMatchesLeft = 2
// 				case int64(SoftInjury):
// 					injuryMatchesLeft = 3
// 				case int64(HardInjury):
// 					injuryMatchesLeft = 7
// 				default:
// 					return fmt.Errorf("out of game type unknown %v", outOfGameType)
// 				}
// 			}
// 		}
// 		if is2ndHalf {
// 			if redCardMatchesLeft > 0 {
// 				redCardMatchesLeft--
// 			}
// 			if injuryMatchesLeft > 0 {
// 				injuryMatchesLeft--
// 			}
// 		}
// 		// log.Infof("encoded skills %v, redCard %v, injuries %v", player.EncodedSkills, player.RedCardMatchesLeft, player.InjuryMatchesLeft)
// 		if err = player.SetRedCard(contracts, redCardMatchesLeft != 0); err != nil {
// 			return err
// 		}
// 		if err = player.SetInjuryWeeks(contracts, injuryMatchesLeft); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (b *Team) updateTeamSkills(
// 	contracts contracts.Contracts,
// 	logs *big.Int,
// 	startTime *big.Int,
// ) error {
// 	trainingPoints, err := contracts.Evolution.GetTrainingPoints(&bind.CallOpts{}, logs)
// 	if err != nil {
// 		return err
// 	}
// 	b.TrainingPoints = trainingPoints.Uint64()

// 	userAssignment, _ := new(big.Int).SetString("1022963800726800053580157736076735226208686447456863237", 10)
// 	newSkills, err := contracts.Evolution.GetTeamEvolvedSkills(
// 		&bind.CallOpts{},
// 		b.Skills(),
// 		userAssignment,
// 		startTime,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	for i, player := range b.Players {
// 		player.sto.EncodedSkills = newSkills[i]
// 	}
// 	// TODO the followign code is old but it contain the code to evolve the name of a change in the generation of a player. Add it in the future.
// 	// for s, state := range newStates {
// 	// 	if state.String() == "0" {
// 	// 		continue
// 	// 	}

// 	// 	playerID, err := b.contracts.Leagues.GetPlayerIdFromSkills(&bind.CallOpts{}, state)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	player, err := storage.PlayerByPlayerId(tx, playerID)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	if player == nil {
// 	// 		return fmt.Errorf("Unexistent playerId %v", playerID)
// 	// 	}
// 	// 	oldGen, err := b.contracts.Assets.GetGeneration(&bind.CallOpts{}, states[s])
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	newGen, err := b.contracts.Assets.GetGeneration(&bind.CallOpts{}, state)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	if newGen.Cmp(oldGen) != 0 {
// 	// 		timezone, countryIdx, _, err := b.contracts.Assets.DecodeTZCountryAndVal(&bind.CallOpts{}, player.PlayerId)
// 	// 		if err != nil {
// 	// 			return err
// 	// 		}
// 	// 		newName, err := b.namesdb.GeneratePlayerFullName(player.PlayerId, uint8(newGen.Uint64()), timezone, countryIdx.Uint64())
// 	// 		if err != nil {
// 	// 			return err
// 	// 		}
// 	// 		player.Name = newName
// 	// 	}
// 	// 	defence, speed, pass, shoot, endurance, _, _, err := utils.DecodeSkills(b.contracts.Assets, state)
// 	// 	player.Defence = defence.Uint64()
// 	// 	player.Speed = speed.Uint64()
// 	// 	player.Pass = pass.Uint64()
// 	// 	player.Shoot = shoot.Uint64()
// 	// 	player.Defence = endurance.Uint64()
// 	// 	player.EncodedSkills = state
// 	// 	err = player.Update(tx)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// }
// 	return nil
// }
