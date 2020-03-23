package engine

import (
	"fmt"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type Training struct {
	storage.Training
}

func NewTraining(sto storage.Training) *Training {
	training := Training{}
	training.Training = sto
	return &training
}

func (b Training) Marshal() string {
	var result string
	result += "{"
	if b.SpecialPlayerShirt == -1 {
		result += fmt.Sprintf("specialPlayerShirt: '%v',", 25)
	} else {
		result += fmt.Sprintf("specialPlayerShirt: '%v',", b.SpecialPlayerShirt)
	}
	result += fmt.Sprintf("goalkeepersDefence: '%v',", b.GoalkeepersDefence)
	result += fmt.Sprintf("goalkeepersSpeed: '%v',", b.GoalkeepersSpeed)
	result += fmt.Sprintf("goalkeepersPass: '%v',", b.GoalkeepersPass)
	result += fmt.Sprintf("goalkeepersShoot: '%v',", b.GoalkeepersShoot)
	result += fmt.Sprintf("goalkeepersEndurance: '%v',", b.GoalkeepersEndurance)
	result += fmt.Sprintf("defendersDefence: '%v',", b.DefendersDefence)
	result += fmt.Sprintf("defendersSpeed: '%v',", b.DefendersSpeed)
	result += fmt.Sprintf("defendersPass: '%v',", b.DefendersPass)
	result += fmt.Sprintf("defendersShoot: '%v',", b.DefendersShoot)
	result += fmt.Sprintf("defendersEndurance: '%v',", b.DefendersEndurance)
	result += fmt.Sprintf("midfieldersDefence: '%v',", b.MidfieldersDefence)
	result += fmt.Sprintf("midfieldersSpeed: '%v',", b.MidfieldersSpeed)
	result += fmt.Sprintf("midfieldersPass: '%v',", b.MidfieldersPass)
	result += fmt.Sprintf("midfieldersShoot: '%v',", b.MidfieldersShoot)
	result += fmt.Sprintf("midfieldersEndurance: '%v',", b.MidfieldersEndurance)
	result += fmt.Sprintf("attackersDefence: '%v',", b.AttackersDefence)
	result += fmt.Sprintf("attackersSpeed: '%v',", b.AttackersSpeed)
	result += fmt.Sprintf("attackersPass: '%v',", b.AttackersPass)
	result += fmt.Sprintf("attackersShoot: '%v',", b.AttackersShoot)
	result += fmt.Sprintf("attackersEndurance: '%v',", b.AttackersEndurance)
	result += fmt.Sprintf("specialPlayerDefence: '%v',", b.SpecialPlayerDefence)
	result += fmt.Sprintf("specialPlayerSpeed: '%v',", b.SpecialPlayerSpeed)
	result += fmt.Sprintf("specialPlayerPass: '%v',", b.SpecialPlayerPass)
	result += fmt.Sprintf("specialPlayerShoot: '%v',", b.SpecialPlayerShoot)
	result += fmt.Sprintf("specialPlayerEndurance: '%v',", b.SpecialPlayerEndurance)
	result += fmt.Sprintf("}")
	return result

}
