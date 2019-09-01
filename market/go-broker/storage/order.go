package storage

import log "github.com/sirupsen/logrus"

type Order struct {
	SellOrder SellOrder
	BuyOrder  BuyOrder
}

func (b *Storage) GetOrders() ([]Order, error) {
	var orders []Order
	rows, err := b.db.Query("SELECT player_sell_orders.playerId, player_sell_orders.price, player_buy_orders.playerId, player_buy_orders.price FROM player_sell_orders INNER JOIN player_buy_orders;")
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		var order Order
		err = rows.Scan(
			&order.SellOrder.PlayerId,
			&order.SellOrder.Price,
			&order.BuyOrder.PlayerId,
			&order.BuyOrder.Price,
		)
		if err != nil {
			return orders, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *Storage) DeleteOrder(playerId uint64) error {
	log.Infof("[DBMS] + create order %v", playerId)
	err := b.DeleteBuyOrder(playerId)
	if err != nil {
		return err
	}
	err = b.DeleteSellOrder(playerId)
	return err
}
