package postgres_test

import (
	"testing"

	"github.com/freeverseio/crypto-soccer/go/storage/postgres"
	"github.com/freeverseio/crypto-soccer/go/storage/storagetest"
	"gotest.tools/assert"
)

func TestTeamStorageServiceInterface(t *testing.T) {
	tx, err := s.Begin()
	assert.NilError(t, err)
	defer tx.Rollback()

	createMinimumUniverse(t, tx)

	service := postgres.NewTeamStorageService(tx)
	storagetest.TestTeamStorageService(t, service)
}
