package storage

import (
	"database/sql"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type SellOrder struct {
	PlayerId   *big.Int
	CurrencyId uint8
	Price      uint64
	Rnd        *big.Int
	ValidUntil *big.Int
	TypeOfTx   int8
	Signature  string
}

func (b *Storage) CreateSellOrder(order SellOrder) error {
	log.Infof("[DBMS] + create sell order %v", order)
	_, err := b.db.Exec("INSERT INTO player_sell_orders (playerId, currencyId, price, rnd, validUntil, typeOfTx, signature) VALUES ($1, $2, $3, $4, $5, $6, $7);",
		order.PlayerId.String(),
		order.CurrencyId,
		order.Price,
		order.Rnd.String(),
		order.ValidUntil.String(),
		order.TypeOfTx,
		order.Signature,
	)
	return err
}

func (b *Storage) GetSellOrders() ([]SellOrder, error) {
	var orders []SellOrder
	rows, err := b.db.Query("SELECT playerId, currencyId, price, rnd, validUntil, typeOfTx, signature FROM player_sell_orders;")
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		var order SellOrder
		var playerId sql.NullString
		var rnd sql.NullString
		var validUntil sql.NullString
		err = rows.Scan(
			&playerId,
			&order.CurrencyId,
			&order.Price,
			&rnd,
			&validUntil,
			&order.TypeOfTx,
			&order.Signature,
		)
		if err != nil {
			return orders, err
		}
		order.PlayerId, _ = new(big.Int).SetString(playerId.String, 10)
		order.Rnd, _ = new(big.Int).SetString(rnd.String, 10)
		order.ValidUntil, _ = new(big.Int).SetString(validUntil.String, 10)
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *Storage) DeleteSellOrder(playerId *big.Int) error {
	log.Infof("[DBMS] Delete sell order %v", playerId)
	_, err := b.db.Exec("DELETE FROM player_sell_orders WHERE (playerId = $1);",
		playerId.String(),
	)
	return err
}
