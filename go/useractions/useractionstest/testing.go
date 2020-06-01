package useractionstest

import (
	"database/sql"
	"math"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"github.com/freeverseio/crypto-soccer/go/useractions"
	"gotest.tools/assert"
)

func TestUserActionsPublishService(t *testing.T, service useractions.UserActionsPublishService) {
	t.Run("TestIpfsPushAndPull", func(t *testing.T) {
		t.Parallel()
		var ua useractions.UserActions
		tactic := storage.Tactic{}
		tactic.TeamID = "ciao"
		ua.Tactics = append(ua.Tactics, tactic)
		id, err := service.Publish(ua)
		assert.NilError(t, err)
		// assert.Equal(t, id, "QmXCYKHSNDCHqzv6W7WDHyW1Zp2YLgt87gmt8tzZYTQtx7")
		training := storage.Training{}
		training.TeamID = "pippo"
		ua.Trainings = append(ua.Trainings, training)
		id, err = service.Publish(ua)
		assert.NilError(t, err)
		// assert.Equal(t, id, "QmbUVhwjGJQzPQQjs5QvJjRZLYuW2jKMKf1RcRiNP71qf2")
		ua2, err := service.Retrive(id)
		assert.NilError(t, err)
		assert.Assert(t, ua2 != nil)
		assert.Assert(t, ua2.Equal(&ua))
	})
}

func TestUserActionsStorageService(t *testing.T, service useractions.UserActionsStorageService) {
	t.Run("TestUserActionsPullFromStorageNoUserActions", func(t *testing.T) {
		timezone := 4
		ua, err := service.UserActionsByTimezone(timezone)
		assert.NilError(t, err)
		assert.Equal(t, len(ua.Tactics), 0)
		assert.Equal(t, len(ua.Trainings), 0)
	})

	t.Run("TestUserActionsPullFromStorage", func(t *testing.T) {
		training := storage.NewTraining()
		training.TeamID = teamID
		tactic := storage.Tactic{}
		tactic.TeamID = teamID
		actions := useractions.UserActions{}
		actions.Trainings = append(actions.Trainings, *training)
		actions.Tactics = append(actions.Tactics, tactic)
		assert.NilError(t, service.Insert(actions))
		ua, err := service.UserActionsByTimezone(int(timezoneIdx))
		assert.NilError(t, err)
		assert.Assert(t, ua.Equal(&actions))
	})
}

const timezoneIdx = uint8(1)
const countryIdx = uint32(0)
const leagueIdx = uint32(0)
const teamID = "1"

func CreateMinimumUniverse(t *testing.T, tx *sql.Tx) {
	timezone := storage.Timezone{timezoneIdx}
	assert.NilError(t, timezone.Insert(tx))

	country := storage.Country{timezone.TimezoneIdx, countryIdx}
	assert.NilError(t, country.Insert(tx))

	league := storage.League{timezone.TimezoneIdx, countryIdx, leagueIdx}
	assert.NilError(t, league.Insert(tx))

	team := storage.NewTeam()
	team.TeamID = teamID
	team.TimezoneIdx = timezone.TimezoneIdx
	team.CountryIdx = countryIdx
	team.LeagueIdx = leagueIdx
	team.Owner = "my team"
	team.RankingPoints = math.MaxUint64
	assert.NilError(t, team.Insert(tx))
}
