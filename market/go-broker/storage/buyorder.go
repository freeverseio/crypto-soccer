package storage

import log "github.com/sirupsen/logrus"

type BuyOrder struct {
	PlayerId uint64
	Price    uint64
	Owner    string
}

func (b *Storage) CreateBuyOrder(order BuyOrder) error {
	log.Infof("[DBMS] + create buy order %v", order)
	_, err := b.db.Exec("INSERT INTO player_buy_orders (playerId, price, owner) VALUES ($1, $2, $3);",
		order.PlayerId,
		order.Price,
		order.Owner,
	)
	return err
}

func (b *Storage) GetBuyOrders() ([]BuyOrder, error) {
	var offers []BuyOrder
	rows, err := b.db.Query("SELECT playerId, price, owner FROM player_buy_orders;")
	if err != nil {
		return offers, err
	}
	defer rows.Close()
	for rows.Next() {
		var offer BuyOrder
		err = rows.Scan(
			&offer.PlayerId,
			&offer.Price,
			&offer.Owner,
		)
		if err != nil {
			return offers, err
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

func (b *Storage) DeleteBuyOrder(playerId uint64) error {
	log.Infof("[DBMS] Delete buy order %v", playerId)
	_, err := b.db.Exec("DELETE FROM player_buy_orders WHERE (playerId == $0);",
		playerId,
	)
	return err
}
