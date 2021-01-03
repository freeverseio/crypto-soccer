package postgres_test

import (
	"database/sql"
	"testing"

	"github.com/freeverseio/crypto-soccer/go/purchasevoider/postgres"
	"gotest.tools/assert"
)

func TestMarketService(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://freeverse:freeverse@crypto-soccer_devcontainer_dockerhost_1:5432/market?sslmode=disable")
	assert.NilError(t, err)
	service := postgres.MarketService{db}
	token, err := service.GetPlayerIdByPurchaseToken("token")
	assert.NilError(t, err)
	assert.Equal(t, token, "")
}
