package engine

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/freeverseio/crypto-soccer/go/contracts"
	"github.com/freeverseio/crypto-soccer/go/storage"
)

type Team struct {
	storage.Team
	Players  [25]*Player
	Training Training
}

func NewTeam() *Team {
	var team Team
	team.Team = *storage.NewTeam()
	for i := range team.Players {
		team.Players[i] = NewPlayer()
	}
	team.Training = *NewTraining(*storage.NewTraining())
	return &team
}

func (b Team) ToStorage(contracts contracts.Contracts, tx *sql.Tx, blockNumber uint64) error {
	for _, player := range b.Players {
		if player.IsNil() {
			continue
		}
		if err := player.Update(tx, blockNumber); err != nil {
			return err
		}
	}
	return b.Update(tx)
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
		b.Players[i].SetSkills(contracts, skills[i])
	}
}

// order: shoot, speed, pass, defence, endurance
func (b Team) EncodeAssignedTrainingPoints(contracts contracts.Contracts) (*big.Int, error) {
	TPperSkill := b.Training.Goalkeepers.ToSlice()
	TPperSkill = append(TPperSkill, b.Training.Defenders.ToSlice()...)
	TPperSkill = append(TPperSkill, b.Training.Midfielders.ToSlice()...)
	TPperSkill = append(TPperSkill, b.Training.Attackers.ToSlice()...)
	TPperSkill = append(TPperSkill, b.Training.SpecialPlayer.ToSlice()...)

	var TPperSkillFixedSize [25]uint16
	copy(TPperSkillFixedSize[:], TPperSkill[:25])
	specialPlayer := uint8(25)
	if b.Training.SpecialPlayerShirt >= 0 && b.Training.SpecialPlayerShirt < 25 {
		specialPlayer = uint8(b.Training.SpecialPlayerShirt)
	}
	encodedTraining, err := contracts.TrainingPoints.EncodeTP(
		&bind.CallOpts{},
		b.TrainingPoints,
		TPperSkillFixedSize,
		specialPlayer,
	)
	if err != nil {
		return nil, err
	}
	return encodedTraining, nil
}

func (b Team) ToJavaScript() string {
	var result string
	result += "{"
	result += fmt.Sprintf("matchLog: '%v',", b.MatchLog)
	result += fmt.Sprintf("teamId: '%v',", b.TeamID)
	result += fmt.Sprintf("tactic: '%v',", b.Tactic)
	result += fmt.Sprintf("trainingPoints: '%v',", b.TrainingPoints)
	result += fmt.Sprintf("training: %v,", b.Training.Marshal())
	result += "players: ["
	for _, player := range b.Players {
		result += fmt.Sprintf("'%v',", player.EncodedSkills)
	}
	result += "],}"
	return result
}
