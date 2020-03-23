package engine

import (
	"database/sql"
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

func (b Team) CalculateAssignedTrainingPoints(contracts contracts.Contracts) (*big.Int, error) {
	TPperSkill := [25]uint16{
		uint16(b.Training.GoalkeepersDefence),
		uint16(b.Training.GoalkeepersSpeed),
		uint16(b.Training.GoalkeepersPass),
		uint16(b.Training.GoalkeepersShoot),
		uint16(b.Training.GoalkeepersEndurance),
		uint16(b.Training.DefendersDefence),
		uint16(b.Training.DefendersSpeed),
		uint16(b.Training.DefendersPass),
		uint16(b.Training.DefendersShoot),
		uint16(b.Training.DefendersEndurance),
		uint16(b.Training.MidfieldersDefence),
		uint16(b.Training.MidfieldersSpeed),
		uint16(b.Training.MidfieldersPass),
		uint16(b.Training.MidfieldersShoot),
		uint16(b.Training.MidfieldersEndurance),
		uint16(b.Training.AttackersDefence),
		uint16(b.Training.AttackersSpeed),
		uint16(b.Training.AttackersPass),
		uint16(b.Training.AttackersShoot),
		uint16(b.Training.AttackersEndurance),
		uint16(b.Training.SpecialPlayerDefence),
		uint16(b.Training.SpecialPlayerSpeed),
		uint16(b.Training.SpecialPlayerPass),
		uint16(b.Training.SpecialPlayerShoot),
		uint16(b.Training.SpecialPlayerEndurance),
	}
	specialPlayer := uint8(25)
	if b.Training.SpecialPlayerShirt >= 0 && b.Training.SpecialPlayerShirt < 25 {
		specialPlayer = uint8(b.Training.SpecialPlayerShirt)
	}
	encodedTraining, err := contracts.TrainingPoints.EncodeTP(
		&bind.CallOpts{},
		b.TrainingPoints,
		TPperSkill,
		specialPlayer,
	)
	if err != nil {
		return nil, err
	}
	return encodedTraining, nil
}
