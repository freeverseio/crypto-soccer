package process

import (
	"errors"
	"math/big"
	"sort"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type OrgMap struct {
	teams []storage.Team
}

func (b *OrgMap) AppendOrgMap(newOrgMap OrgMap) error {
	for _, team := range newOrgMap.teams {
		b.Append(team)
	}
	return nil
}

func (b *OrgMap) Append(team storage.Team) error {
	if teamID, _ := new(big.Int).SetString(team.TeamID, 10); teamID == nil {
		return errors.New("invalid TeamID")
	}

	for _, t := range b.teams {
		if t.TeamID == team.TeamID {
			return errors.New("team already in OrgMap")
		}
	}

	b.teams = append(b.teams, team)
	return nil
}

func (b *OrgMap) Sort() {
	sort.Slice(b.teams[:], func(i, j int) bool {
		if b.teams[i].RankingPoints == b.teams[j].RankingPoints {
			teamID0, _ := new(big.Int).SetString(b.teams[i].TeamID, 10)
			teamID1, _ := new(big.Int).SetString(b.teams[j].TeamID, 10)
			return teamID0.Cmp(teamID1) == -1
		}
		return b.teams[i].RankingPoints > b.teams[j].RankingPoints
	})
}

func (b OrgMap) Size() int {
	return len(b.teams)
}

func (b OrgMap) At(i int) storage.Team {
	return b.teams[i]
}
