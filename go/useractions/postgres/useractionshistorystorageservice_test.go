package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/useractions/postgres"
	"github.com/freeverseio/crypto-soccer/go/useractions/useractionstest"
	"gotest.tools/assert"
)

func TestUserActionsHistoryStorageInterface(t *testing.T) {
	t.Parallel()
	tx, err := db.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()
	useractionstest.CreateMinimumUniverse(t, tx)
	blockNumber := uint64(23123)
	service := postgres.NewUserActionsHistoryStorageService(tx, blockNumber)
	useractionstest.TestUserActionsStorageService(t, service)
}
