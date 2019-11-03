package storage

import (
	"database/sql"
	"math/big"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type BidState string

const (
	BID_ACCEPTED      BidState = "ACCEPTED"
	BID_REFUSED       BidState = "REFUSED"
	BID_PAYING        BidState = "PAYING"
	BID_PAID          BidState = "PAID"
	BID_FAILED_TO_PAY BidState = "FAILED_TO_PAY"
)

type Bid struct {
	Auction         uuid.UUID
	ExtraPrice      int64
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

func (b *Storage) UpdateBidState(auction uuid.UUID, extra_price int64, state BidState) error {
	_, err := b.db.Exec("UPDATE bids SET state=$1 WHERE auction=$2 AND extra_price=$3;", state, auction, extra_price)
	return err
}

func (b *Storage) UpdateBidPaymentUrl(auction uuid.UUID, extra_price int64, url string) error {
	_, err := b.db.Exec("UPDATE bids SET payment_url=$1 WHERE auction=$2 AND extra_price=$3;", url, auction, extra_price)
	return err
}

func (b *Storage) GetBidsOfAuction(auctionUUID uuid.UUID) ([]Bid, error) {
	var bids []Bid
	rows, err := b.db.Query("SELECT extra_price, rnd, team_id, signature, state FROM bids WHERE auction=$1;", auctionUUID)
	if err != nil {
		return bids, err
	}
	defer rows.Close()
	for rows.Next() {
		var bid Bid
		var teamID sql.NullString
		err = rows.Scan(
			&bid.ExtraPrice,
			&bid.Rnd,
			&teamID,
			&bid.Signature,
			&bid.State,
		)
		if err != nil {
			return bids, err
		}
		bid.Auction = auctionUUID
		bid.TeamID, _ = new(big.Int).SetString(teamID.String, 10)
		bids = append(bids, bid)
	}
	return bids, nil
}
