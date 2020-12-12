package postgres_test

import (
	"database/sql"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/purchasevoider/postgres"
	"gotest.tools/assert"
)

func TestUniverseService(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable")
	assert.NilError(t, err)
	service := postgres.UniverseService{db}
	assert.NilError(t, service.MarkForDeletion("id"))
}
