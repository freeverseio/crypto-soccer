package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type MarketService struct {
	DB *sql.DB
}

func (b *MarketService) GetPlayerIdByPurchaseToken(token string) (string, error) {
	rows, err := b.DB.Query("SELECT player_id FROM playstore_orders WHERE purchase_token=$1", token)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", err
	}

	var id string
	if err := rows.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}
