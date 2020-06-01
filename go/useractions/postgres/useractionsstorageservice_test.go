package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/useractions/postgres"
	"github.com/freeverseio/crypto-soccer/go/useractions/useractionstest"
	"gotest.tools/assert"
)

func TestUserActionsStorageInterface(t *testing.T) {
	t.Parallel()
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	useractionstest.CreateMinimumUniverse(t, tx)
	service := postgres.NewUserActionsStorageService(tx)
	useractionstest.TestUserActionsStorageService(t, service)
}
