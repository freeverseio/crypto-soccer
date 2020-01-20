package engine

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
)

const NoOutOfGamePlayer = uint8(14)
const RedCard = uint8(3)
const SoftInjury = uint8(1)
const HardInjury = uint8(2)

type Team struct {
	TeamID         *big.Int
	Players        [25]*Player
	tactic         *big.Int
	TrainingPoints uint64
}

func NewTeam(
	contracts *contracts.Contracts,
) (*Team, error) {
	var team Team
	team.TeamID = big.NewInt(0)
	for i := range team.Players {
		team.Players[i] = NewNullPlayer()
	}
	var err error
	if team.tactic, err = DefaultTactic(contracts); err != nil {
		return nil, err
	}
	return &team, nil
}

func (b Team) DumpState() string {
	var state string
	state += fmt.Sprintf("TeamId: %v\n", b.TeamID)
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

func DefaultTactic(contracts *contracts.Contracts) (*big.Int, error) {
	substitutions := [3]uint8{11, 11, 11}
	substitutionsMinute := [3]uint8{2, 3, 4}
	formation := [14]uint8{0, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 25, 25, 25}
	extraAttack := [10]bool{false, false, false, false, false, false, false, false, false, false}
	tacticID := uint8(1)
	tactic, err := contracts.Engine.EncodeTactics(
		&bind.CallOpts{},
		substitutions,
		substitutionsMinute,
		formation,
		extraAttack,
		tacticID,
	)
	if err != nil {
		return nil, err
	}
	return tactic, nil
}

func (b *Team) updateTeamState(
	contracts contracts.Contracts,
	is2ndHalf bool,
	matchLog *big.Int,
) error {
	decodedTactic, err := contracts.Leagues.DecodeTactics(&bind.CallOpts{}, b.tactic)
	if err != nil {
		return err
	}
	outOfGamePlayer, err := contracts.Engineprecomp.GetOutOfGamePlayer(&bind.CallOpts{}, matchLog, is2ndHalf)
	if err != nil {
		return err
	}
	outOfGameType, err := contracts.Engineprecomp.GetOutOfGameType(&bind.CallOpts{}, matchLog, is2ndHalf)
	if err != nil {
		return err
	}
	for shirt, player := range b.Players {
		wasAligned, err := contracts.Engine.WasPlayerAlignedEndOfLastHalf(
			&bind.CallOpts{},
			uint8(shirt),
			b.tactic,
			matchLog,
		)
		if err != nil {
			return err
		}
		player.skills, err = contracts.Evolution.SetAlignedEndOfLastHalf(
			&bind.CallOpts{},
			player.skills,
			wasAligned,
		)
		if err != nil {
			return err
		}
		var redCardMatchesLeft uint8
		var injuryMatchesLeft uint8
		if outOfGamePlayer.Int64() != int64(NoOutOfGamePlayer) {
			if outOfGamePlayer.Int64() < 0 || int(outOfGamePlayer.Int64()) >= len(decodedTactic.Lineup) {
				return fmt.Errorf("out of game player unknown %v, tactics %v, matchlog %v", outOfGamePlayer.Int64(), b.tactic, matchLog)
			}
			if uint8(shirt) == decodedTactic.Lineup[outOfGamePlayer.Int64()] {
				switch outOfGameType.Int64() {
				case int64(RedCard):
					redCardMatchesLeft = 2
				case int64(SoftInjury):
					injuryMatchesLeft = 3
				case int64(HardInjury):
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
		if player.skills, err = contracts.Evolution.SetRedCardLastGame(&bind.CallOpts{}, player.skills, redCardMatchesLeft != 0); err != nil {
			return err
		}
		if player.skills, err = contracts.Evolution.SetInjuryWeeksLeft(&bind.CallOpts{}, player.skills, injuryMatchesLeft); err != nil {
			return err
		}

	}
	return nil
}

func (b *Team) SetTrainingPoints(contracts *contracts.Contracts, matchLog *big.Int) error {
	points, err := contracts.Evolution.GetTrainingPoints(&bind.CallOpts{}, matchLog)
	if err != nil {
		return err
	}
	b.TrainingPoints = points.Uint64()
	return nil
}
