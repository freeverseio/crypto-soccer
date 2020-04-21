package process_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/synchronizer/process"
	"gotest.tools/assert"
)

func TestOrgMapAppend(t *testing.T) {
	om := process.OrgMap{}
	assert.Equal(t, om.Size(), 0)
	assert.NilError(t, om.Append(*storage.NewTeam()))
	assert.Equal(t, om.Size(), 1)
}

func TestOrgMapAddTwiceSameElement(t *testing.T) {
	team := storage.NewTeam()
	om := process.OrgMap{}
	assert.NilError(t, om.Append(*team))
	assert.Error(t, om.Append(*team), "team already in OrgMap")
}

func TestOrgMapAddWithInvalidTeamID(t *testing.T) {
	team := storage.NewTeam()
	team.TeamID = "ciao"
	om := process.OrgMap{}
	assert.Error(t, om.Append(*team), "invalid TeamID")
}

func TestOrgMapAddTeamWithSameRankingPoints(t *testing.T) {
	team := storage.NewTeam()
	team.TeamID = "1"
	om := process.OrgMap{}
	assert.NilError(t, om.Append(*team))
	team.TeamID = "2"
	assert.NilError(t, om.Append(*team))
}

func TestOrgMapSort(t *testing.T) {
	team0 := storage.NewTeam()
	team0.TeamID = "3"
	team0.RankingPoints = 0
	team1 := storage.NewTeam()
	team1.TeamID = "5"
	team1.RankingPoints = 0
	team2 := storage.NewTeam()
	team2.TeamID = "7"
	team2.RankingPoints = 2
	om := process.OrgMap{}
	assert.NilError(t, om.Append(*team0))
	assert.NilError(t, om.Append(*team1))
	assert.NilError(t, om.Append(*team2))
	om.Sort()
	assert.Equal(t, om.At(0), *team2)
	assert.Equal(t, om.At(1), *team0)
	assert.Equal(t, om.At(2), *team1)

}
