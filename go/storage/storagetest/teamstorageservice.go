package storagetest

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestTeamStorageService(t *testing.T, service storage.TeamStorageService) {
	t.Run("insert a team", func(t *testing.T) {
		team := storage.NewTeam()
		team.TeamID = "2"
		team.TimezoneIdx = 1
		team.CountryIdx = 0
		team.LeagueIdx = 0
		team.Name = "pippo"
		assert.NilError(t, service.Insert(*team))
	})
	t.Run("update name of unexistent team", func(t *testing.T) {
		assert.Error(t, service.UpdateName("", ""), "unexistent team")
	})
	t.Run("update name of team", func(t *testing.T) {
		team := storage.NewTeam()
		team.TeamID = "3"
		team.TimezoneIdx = 1
		team.CountryIdx = 0
		team.LeagueIdx = 0
		team.Name = "pippo"
		assert.NilError(t, service.Insert(*team))

		service.UpdateName(team.TeamID, "pippo2")
		resultTeam, err := service.Team(team.TeamID)
		assert.NilError(t, err)
		assert.Equal(t, resultTeam.Name, "pippo2")
	})
}
