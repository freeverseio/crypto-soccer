package process

import (
	"errors"
	"math/big"
	"sort"

	"github.com/freeverseio/crypto-soccer/go/storage"
)

type OrgMap struct {
	Teams []storage.Team
}

func (b *OrgMap) AppendOrgMap(newOrgMap OrgMap) error {
	for _, team := range newOrgMap.Teams {
		b.Append(team)
	}
	return nil
}

func (b *OrgMap) Append(team storage.Team) error {
	if teamID, _ := new(big.Int).SetString(team.TeamID, 10); teamID == nil {
		return errors.New("invalid TeamID")
	}

	for _, t := range b.Teams {
		if t.TeamID == team.TeamID {
			return errors.New("team already in OrgMap")
		}
	}

	b.Teams = append(b.Teams, team)
	return nil
}

func (b *OrgMap) Sort() {
	sort.Slice(b.Teams[:], func(i, j int) bool {
		if b.Teams[i].RankingPoints == b.Teams[j].RankingPoints {
			teamID0, _ := new(big.Int).SetString(b.Teams[i].TeamID, 10)
			teamID1, _ := new(big.Int).SetString(b.Teams[j].TeamID, 10)
			return teamID0.Cmp(teamID1) == -1
		}
		return b.Teams[i].RankingPoints > b.Teams[j].RankingPoints
	})
}

func (b OrgMap) Size() int {
	return len(b.Teams)
}

func (b OrgMap) At(i int) storage.Team {
	return b.Teams[i]
}
