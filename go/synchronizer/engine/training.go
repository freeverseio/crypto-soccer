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
	result += fmt.Sprintf("goalkeepersDefence: '%v',", b.Goalkeepers.Defence)
	result += fmt.Sprintf("goalkeepersSpeed: '%v',", b.Goalkeepers.Speed)
	result += fmt.Sprintf("goalkeepersPass: '%v',", b.Goalkeepers.Pass)
	result += fmt.Sprintf("goalkeepersShoot: '%v',", b.Goalkeepers.Shoot)
	result += fmt.Sprintf("goalkeepersEndurance: '%v',", b.Goalkeepers.Endurance)
	result += fmt.Sprintf("defendersDefence: '%v',", b.Defenders.Defence)
	result += fmt.Sprintf("defendersSpeed: '%v',", b.Defenders.Speed)
	result += fmt.Sprintf("defendersPass: '%v',", b.Defenders.Pass)
	result += fmt.Sprintf("defendersShoot: '%v',", b.Defenders.Shoot)
	result += fmt.Sprintf("defendersEndurance: '%v',", b.Defenders.Endurance)
	result += fmt.Sprintf("midfieldersDefence: '%v',", b.Midfielders.Defence)
	result += fmt.Sprintf("midfieldersSpeed: '%v',", b.Midfielders.Speed)
	result += fmt.Sprintf("midfieldersPass: '%v',", b.Midfielders.Pass)
	result += fmt.Sprintf("midfieldersShoot: '%v',", b.Midfielders.Shoot)
	result += fmt.Sprintf("midfieldersEndurance: '%v',", b.Midfielders.Endurance)
	result += fmt.Sprintf("attackersDefence: '%v',", b.Attackers.Defence)
	result += fmt.Sprintf("attackersSpeed: '%v',", b.Attackers.Speed)
	result += fmt.Sprintf("attackersPass: '%v',", b.Attackers.Pass)
	result += fmt.Sprintf("attackersShoot: '%v',", b.Attackers.Shoot)
	result += fmt.Sprintf("attackersEndurance: '%v',", b.Attackers.Endurance)
	result += fmt.Sprintf("specialPlayerDefence: '%v',", b.SpecialPlayer.Defence)
	result += fmt.Sprintf("specialPlayerSpeed: '%v',", b.SpecialPlayer.Speed)
	result += fmt.Sprintf("specialPlayerPass: '%v',", b.SpecialPlayer.Pass)
	result += fmt.Sprintf("specialPlayerShoot: '%v',", b.SpecialPlayer.Shoot)
	result += fmt.Sprintf("specialPlayerEndurance: '%v',", b.SpecialPlayer.Endurance)
	result += fmt.Sprintf("}")
	return result

}
