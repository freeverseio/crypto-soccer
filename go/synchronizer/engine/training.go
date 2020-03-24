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
	result += fmt.Sprintf("Goalkeepers: {")
	result += fmt.Sprintf("Shoot: '%v',", b.Goalkeepers.Shoot)
	result += fmt.Sprintf("Speed: '%v',", b.Goalkeepers.Speed)
	result += fmt.Sprintf("Pass: '%v',", b.Goalkeepers.Pass)
	result += fmt.Sprintf("Defence: '%v',", b.Goalkeepers.Defence)
	result += fmt.Sprintf("Endurance: '%v',", b.Goalkeepers.Endurance)
	result += fmt.Sprintf("}, ")
	result += fmt.Sprintf("Defenders: {")
	result += fmt.Sprintf("Shoot: '%v',", b.Defenders.Shoot)
	result += fmt.Sprintf("Speed: '%v',", b.Defenders.Speed)
	result += fmt.Sprintf("Pass: '%v',", b.Defenders.Pass)
	result += fmt.Sprintf("Defence: '%v',", b.Defenders.Defence)
	result += fmt.Sprintf("Endurance: '%v',", b.Defenders.Endurance)
	result += fmt.Sprintf("}, ")
	result += fmt.Sprintf("Midfielders: {")
	result += fmt.Sprintf("Shoot: '%v',", b.Midfielders.Shoot)
	result += fmt.Sprintf("Speed: '%v',", b.Midfielders.Speed)
	result += fmt.Sprintf("Pass: '%v',", b.Midfielders.Pass)
	result += fmt.Sprintf("Defence: '%v',", b.Midfielders.Defence)
	result += fmt.Sprintf("Endurance: '%v',", b.Midfielders.Endurance)
	result += fmt.Sprintf("}, ")
	result += fmt.Sprintf("Attackers: {")
	result += fmt.Sprintf("Shoot: '%v',", b.Attackers.Shoot)
	result += fmt.Sprintf("Speed: '%v',", b.Attackers.Speed)
	result += fmt.Sprintf("Pass: '%v',", b.Attackers.Pass)
	result += fmt.Sprintf("Defence: '%v',", b.Attackers.Defence)
	result += fmt.Sprintf("Endurance: '%v',", b.Attackers.Endurance)
	result += fmt.Sprintf("}, ")
	result += fmt.Sprintf("SpecialPlayer: {")
	result += fmt.Sprintf("Shoot: '%v',", b.SpecialPlayer.Shoot)
	result += fmt.Sprintf("Speed: '%v',", b.SpecialPlayer.Speed)
	result += fmt.Sprintf("Pass: '%v',", b.SpecialPlayer.Pass)
	result += fmt.Sprintf("Defence: '%v',", b.SpecialPlayer.Defence)
	result += fmt.Sprintf("Endurance: '%v',", b.SpecialPlayer.Endurance)
	result += fmt.Sprintf("}, ")
	result += fmt.Sprintf("}")
	return result
}
