package storage

import (
	"database/sql"
	"math/big"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Bid struct {
	Auction         uuid.UUID
	ExtraPrice      float32
	Rnd             int64
	TeamID          *big.Int
	Is2StartAuction bool
	Signature       string
	State           BidState
}

func (b *Storage) CreateBid(bid Bid) error {
	log.Infof("[DBMS] + create Bid %v", bid)
	_, err := b.db.Exec("INSERT INTO bids (auction, extra_price, rnd, team_id, signature, state) VALUES ($1, $2, $3, $4, $5, $6);",
		bid.Auction,
		bid.ExtraPrice,
		bid.Rnd,
		bid.TeamID.String(),
		bid.Signature,
		bid.State,
	)
	return err
}

func (b *Storage) GetBids() ([]Bid, error) {
	var bids []Bid
	rows, err := b.db.Query("SELECT auction, extra_price, rnd, team_id, signature, state FROM bids;")
	if err != nil {
		return bids, err
	}
	defer rows.Close()
	for rows.Next() {
		var bid Bid
		var teamID sql.NullString
		err = rows.Scan(
			&bid.Auction,
			&bid.ExtraPrice,
			&bid.Rnd,
			&teamID,
			&bid.Signature,
			&bid.State,
		)
		if err != nil {
			return bids, err
		}
		bid.TeamID, _ = new(big.Int).SetString(teamID.String, 10)
		bids = append(bids, bid)
	}
	return bids, nil
}
