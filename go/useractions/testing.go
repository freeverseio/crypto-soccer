package useractions

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage"
	"gotest.tools/assert"
)

func TestUserActionsPublishService(t *testing.T, service UserActionsPublishService) {
	t.Run("TestIpfsPushAndPull", func(t *testing.T) {
		t.Parallel()
		var ua UserActions
		tactic := storage.Tactic{}
		tactic.TeamID = "ciao"
		ua.Tactics = append(ua.Tactics, tactic)
		cif, err := service.Publish(ua)
		assert.NilError(t, err)
		assert.Equal(t, cif, "QmXCYKHSNDCHqzv6W7WDHyW1Zp2YLgt87gmt8tzZYTQtx7")
		training := storage.Training{}
		training.TeamID = "pippo"
		ua.Trainings = append(ua.Trainings, training)
		cif, err = service.Publish(ua)
		assert.NilError(t, err)
		assert.Equal(t, cif, "QmbUVhwjGJQzPQQjs5QvJjRZLYuW2jKMKf1RcRiNP71qf2")
		ua2, err := service.Retrive(cif)
		assert.NilError(t, err)
		assert.Assert(t, ua2.Equal(&ua))
	})
}
