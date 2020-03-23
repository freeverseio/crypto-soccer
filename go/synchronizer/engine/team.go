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

func SerializeTrainingPerFieldPos(tr storage.TrainingPerFieldPos) []uint16 {
	return []uint16{
		uint16(tr.Shoot),
		uint16(tr.Speed),
		uint16(tr.Pass),
		uint16(tr.Defence),
		uint16(tr.Endurance),
	}
}

func checkTrainingPerFieldPos(availableTPs int, tr storage.TrainingPerFieldPos) [2]bool {
	errTooManyOneSkill := (tr.Shoot + tr.Speed + tr.Pass + tr.Defence + tr.Endurance) > availableTPs
	errTooMany := (100*tr.Shoot > 60*availableTPs) ||
		(100*tr.Speed > 60*availableTPs) ||
		(100*tr.Pass > 60*availableTPs) ||
		(100*tr.Defence > 60*availableTPs) ||
		(100*tr.Endurance > 60*availableTPs)
	return [2]bool{errTooMany, errTooManyOneSkill}
}

func checkTraining(availableTPs int, tr storage.Training) [3]bool {
	errTooMany := false
	errTooManyOneSkill := false
	err := checkTrainingPerFieldPos(availableTPs, tr.Goalkeepers)
	errTooMany = errTooMany || err[0]
	errTooManyOneSkill = errTooMany || err[1]
	err = checkTrainingPerFieldPos(availableTPs, tr.Defenders)
	errTooMany = errTooMany || err[0]
	errTooManyOneSkill = errTooMany || err[1]
	err = checkTrainingPerFieldPos(availableTPs, tr.Midfielders)
	errTooMany = errTooMany || err[0]
	errTooManyOneSkill = errTooMany || err[1]
	err = checkTrainingPerFieldPos(availableTPs, tr.Attackers)
	errTooMany = errTooMany || err[0]
	errTooManyOneSkill = errTooMany || err[1]
	// Special Player has extra 10% points, calculated in this precise integer-division manner:
	availableTPs = (availableTPs * 11) / 10
	err = checkTrainingPerFieldPos(availableTPs, tr.Attackers)
	errSpecialPlayer := err[0] || err[1]
	return [3]bool{errTooMany, errTooManyOneSkill, errSpecialPlayer}
}

// order: shoot, speed, pass, defence, endurance
func (b Team) EncodeAssignedTrainingPoints(contracts contracts.Contracts) (*big.Int, error) {
	TPperSkill := SerializeTrainingPerFieldPos(b.Training.Goalkeepers)
	TPperSkill = append(TPperSkill, SerializeTrainingPerFieldPos(b.Training.Defenders)...)
	TPperSkill = append(TPperSkill, SerializeTrainingPerFieldPos(b.Training.Midfielders)...)
	TPperSkill = append(TPperSkill, SerializeTrainingPerFieldPos(b.Training.Attackers)...)
	TPperSkill = append(TPperSkill, SerializeTrainingPerFieldPos(b.Training.SpecialPlayer)...)

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
