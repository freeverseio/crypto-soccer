package engine

import "github.com/freeverseio/crypto-soccer/go/storage"

type Training struct {
	storage.Training
}

func NewTraining() *Training {
	training := Training{}
	training.Training = *(storage.NewTraining())
	return &training
}

func (b Training) ToString() string {
	var result string
	// result += fmt.Sprintf("seed = '0x%v';", hex.EncodeToString(b.Seed[:]))
	// result += fmt.Sprintf("startTime = '%v';", b.StartTime)
	// result += fmt.Sprintf("matchLog0 = '%v';", b.HomeTeam.MatchLog)
	// result += fmt.Sprintf("teamId0 = '%v';", b.HomeTeam.TeamID)
	// result += fmt.Sprintf("tactic0 = '%v';", b.HomeTeam.Tactic)
	// // result += fmt.Sprintf("assignedTP0 = '%v';", b.HomeTeam.AssignedTP)
	// result += "players0 = ["
	// for _, player := range b.HomeTeam.Players {
	// 	result += fmt.Sprintf("'%v',", player.EncodedSkills)
	// }
	// result += "];"
	// result += fmt.Sprintf("matchLog1 = '%v';", b.VisitorTeam.MatchLog)
	// result += fmt.Sprintf("teamId1 = '%v';", b.VisitorTeam.TeamID)
	// result += fmt.Sprintf("tactic1 = '%v';", b.VisitorTeam.Tactic)
	// // result += fmt.Sprintf("assignedTP1 = '%v';", b.VisitorTeam.AssignedTP)
	// result += "players1 = ["
	// for _, player := range b.VisitorTeam.Players {
	// 	result += fmt.Sprintf("'%v',", player.EncodedSkills)
	// }
	// result += "];"
	return result
}
