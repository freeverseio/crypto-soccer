package storage

import (
	"math/big"

	log "github.com/sirupsen/logrus"
)

type SellOrder struct {
	PlayerId   uint64
	Price      uint64
	Rnd        *big.Int
	ValidUntil *big.Int
	TypeOfTx   int8
}

func (b *Storage) CreateSellOrder(order SellOrder) error {
	log.Infof("[DBMS] + create sell order %v", order)
	_, err := b.db.Exec("INSERT INTO player_sell_orders (playerId, price, rnd, validUntil, typeOfTx) VALUES ($1, $2, $3, $4, $5);",
		order.PlayerId,
		order.Price,
		order.Rnd.String(),
		order.ValidUntil.String(),
		order.TypeOfTx,
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
		err = rows.Scan(
			&order.PlayerId,
			&order.Price,
		)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *Storage) DeleteSellOrder(playerId uint64) error {
	log.Infof("[DBMS] Delete sell order %v", playerId)
	_, err := b.db.Exec("DELETE FROM player_sell_orders WHERE (playerId = $1);",
		playerId,
	)
	return err
}
