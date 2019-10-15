package storage

import (
	"database/sql"
	"math/big"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Bet struct {
	Auction         uuid.UUID
	ExtraPrice      float32
	Rnd             int64
	TeamID          *big.Int
	Is2StartAuction bool
	Signature       string
}

func (b *Storage) CreateBet(order Bet) error {
	log.Infof("[DBMS] + create bet %v", order)
	_, err := b.db.Exec("INSERT INTO bets (auction, extra_price, rnd, team_id, signature) VALUES ($1, $2, $3, $4, $5);",
		order.Auction,
		order.ExtraPrice,
		order.Rnd,
		order.TeamID.String(),
		order.Signature,
	)
	return err
}

func (b *Storage) GetBets() ([]Bet, error) {
	var offers []Bet
	rows, err := b.db.Query("SELECT auction, extra_price, rnd, team_id, signature FROM bets;")
	if err != nil {
		return offers, err
	}
	defer rows.Close()
	for rows.Next() {
		var offer Bet
		var teamID sql.NullString
		err = rows.Scan(
			&offer.Auction,
			&offer.ExtraPrice,
			&offer.Rnd,
			&teamID,
			&offer.Signature,
		)
		if err != nil {
			return offers, err
		}
		offer.TeamID, _ = new(big.Int).SetString(teamID.String, 10)
		offers = append(offers, offer)
	}
	return offers, nil
}
