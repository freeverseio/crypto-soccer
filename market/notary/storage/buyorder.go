package storage

import (
	"database/sql"
	"math/big"

	log "github.com/sirupsen/logrus"
)

type BuyOrder struct {
	PlayerID        *big.Int
	ExtraPrice      float32
	Rnd             int64
	TeamID          *big.Int
	Is2StartAuction bool
	Signature       string
}

func (b *Storage) CreateBuyOrder(order BuyOrder) error {
	log.Infof("[DBMS] + create buy order %v", order)
	_, err := b.db.Exec("INSERT INTO player_buy_orders (player_id, extra_price, rnd, team_id, is_2_start_auction, signature) VALUES ($1, $2, $3, $4, $5, $6);",
		order.PlayerID.String(),
		order.ExtraPrice,
		order.Rnd,
		order.TeamID.String(),
		order.Is2StartAuction,
		order.Signature,
	)
	return err
}

func (b *Storage) GetBuyOrders() ([]BuyOrder, error) {
	var offers []BuyOrder
	rows, err := b.db.Query("SELECT player_id, extra_price, rnd, team_id, is_2_start_auction, signature FROM player_buy_orders;")
	if err != nil {
		return offers, err
	}
	defer rows.Close()
	for rows.Next() {
		var offer BuyOrder
		var playerID sql.NullString
		var teamID sql.NullString
		err = rows.Scan(
			&playerID,
			&offer.ExtraPrice,
			&offer.Rnd,
			&teamID,
			&offer.Is2StartAuction,
			&offer.Signature,
		)
		if err != nil {
			return offers, err
		}
		offer.PlayerID, _ = new(big.Int).SetString(playerID.String, 10)
		offer.TeamID, _ = new(big.Int).SetString(teamID.String, 10)
		offers = append(offers, offer)
	}
	return offers, nil
}

func (b *Storage) DeleteBuyOrder(playerId *big.Int) error {
	log.Infof("[DBMS] Delete buy order %v", playerId)
	_, err := b.db.Exec("DELETE FROM player_buy_orders WHERE (player_id = $1);",
		playerId.String(),
	)
	return err
}
