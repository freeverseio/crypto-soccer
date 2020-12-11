package postgres

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type MarketService struct {
	DB *sql.DB
}

func (b *MarketService) GetPlayerIdByPurchaseToken(token string) (string, error) {
	return "", errors.New("not implemented")
}
