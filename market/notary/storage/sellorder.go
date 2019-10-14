package storage

import (
	"database/sql"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type SellOrder struct {
	PlayerID   *big.Int
	CurrencyID uint8
	Price      *big.Int
	Rnd        *big.Int
	ValidUntil *big.Int
	Signature  string
}

func (b *Storage) CreateSellOrder(order SellOrder) error {
	log.Infof("[DBMS] + create sell order %v", order)
	_, err := b.db.Exec("INSERT INTO player_sell_orders (player_id, currency_id, price, rnd, valid_until, signature) VALUES ($1, $2, $3, $4, $5, $6);",
		order.PlayerID.String(),
		order.CurrencyID,
		order.Price.String(),
		order.Rnd.String(),
		order.ValidUntil.String(),
		order.Signature,
	)
	return err
}

func (b *Storage) GetSellOrders() ([]SellOrder, error) {
	var orders []SellOrder
	rows, err := b.db.Query("SELECT player_id, currency_id, price, rnd, valid_until, signature FROM player_sell_orders;")
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		var order SellOrder
		var playerID sql.NullString
		var price sql.NullString
		var rnd sql.NullString
		var validUntil sql.NullString
		err = rows.Scan(
			&playerID,
			&order.CurrencyID,
			&price,
			&rnd,
			&validUntil,
			&order.Signature,
		)
		if err != nil {
			return orders, err
		}
		order.PlayerID, _ = new(big.Int).SetString(playerID.String, 10)
		order.Price, _ = new(big.Int).SetString(price.String, 10)
		order.Rnd, _ = new(big.Int).SetString(rnd.String, 10)
		order.ValidUntil, _ = new(big.Int).SetString(validUntil.String, 10)
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *Storage) DeleteSellOrder(playerId *big.Int) error {
	log.Infof("[DBMS] Delete sell order %v", playerId)
	_, err := b.db.Exec("DELETE FROM player_sell_orders WHERE (player_id = $1);",
		playerId.String(),
	)
	return err
}
