package storage

import (
	"database/sql"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type BuyOrder struct {
	PlayerId  *big.Int
	TeamId    *big.Int
	Signature string
}

func (b *Storage) CreateBuyOrder(order BuyOrder) error {
	log.Infof("[DBMS] + create buy order %v", order)
	_, err := b.db.Exec("INSERT INTO player_buy_orders (playerId, teamId, signature) VALUES ($1, $2, $3);",
		order.PlayerId.String(),
		order.TeamId.String(),
		order.Signature,
	)
	return err
}

func (b *Storage) GetBuyOrders() ([]BuyOrder, error) {
	var offers []BuyOrder
	rows, err := b.db.Query("SELECT playerId, teamId, signature FROM player_buy_orders;")
	if err != nil {
		return offers, err
	}
	defer rows.Close()
	for rows.Next() {
		var offer BuyOrder
		var playerId sql.NullString
		var teamId sql.NullString
		err = rows.Scan(
			&playerId,
			&teamId,
			&offer.Signature,
		)
		if err != nil {
			return offers, err
		}
		offer.PlayerId, _ = new(big.Int).SetString(playerId.String, 10)
		offer.TeamId, _ = new(big.Int).SetString(teamId.String, 10)
		offers = append(offers, offer)
	}
	return offers, nil
}

func (b *Storage) DeleteBuyOrder(playerId *big.Int) error {
	log.Infof("[DBMS] Delete buy order %v", playerId)
	_, err := b.db.Exec("DELETE FROM player_buy_orders WHERE (playerId = $1);",
		playerId.String(),
	)
	return err
}
