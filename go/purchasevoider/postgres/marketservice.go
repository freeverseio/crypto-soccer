package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/purchasevoider"
	_ "github.com/lib/pq"
)

type MarketService struct {
	db *sql.DB
}

func NewMarketService(db *sql.DB) purchasevoider.MarketService {
	return &MarketService{
		db: db,
	}
}
func (b *MarketService) GetPlayerIdByPurchaseToken(token string) (string, error) {
	return "", nil
}
