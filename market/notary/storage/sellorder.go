package storage

import (
	"database/sql"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type SellOrder struct {
	PlayerId   *big.Int
	Price      uint64
	Rnd        *big.Int
	ValidUntil *big.Int
	TypeOfTx   int8
	Signature  string
}

func (b *Storage) CreateSellOrder(order SellOrder) error {
	log.Infof("[DBMS] + create sell order %v", order)
	_, err := b.db.Exec("INSERT INTO player_sell_orders (playerId, price, rnd, validUntil, typeOfTx, signature) VALUES ($1, $2, $3, $4, $5, $6);",
		order.PlayerId.String(),
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
	rows, err := b.db.Query("SELECT playerId, price FROM player_sell_orders;")
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		var order SellOrder
		var playerId sql.NullString
		err = rows.Scan(
			&playerId,
			&order.Price,
		)
		if err != nil {
			return orders, err
		}
		order.PlayerId, _ = new(big.Int).SetString(playerId.String, 10)
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
