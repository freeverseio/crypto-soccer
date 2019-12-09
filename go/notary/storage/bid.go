package storage

import (
	"database/sql"
	"errors"
	"math/big"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type BidState string

const (
	BIDACCEPTED BidState = "ACCEPTED"
	BIDREFUSED  BidState = "REFUSED"
	BIDPAYING   BidState = "PAYING"
	BIDPAID     BidState = "PAID"
	BIDFAILED   BidState = "FAILED"
)

type Bid struct {
	Auction         uuid.UUID
	ExtraPrice      int64
	Rnd             int64
	TeamID          *big.Int
	Is2StartAuction bool
	Signature       string
	State           BidState
	StateExtra      string
	PaymentID       string
	PaymentURL      string
	PaymentDeadline *big.Int
}

func (b *Storage) CreateBid(bid Bid) error {
	log.Infof("[DBMS] + create Bid %v", bid)
	_, err := b.tx.Exec("INSERT INTO bids (auction, extra_price, rnd, team_id, signature, state) VALUES ($1, $2, $3, $4, $5, $6);",
		bid.Auction,
		bid.ExtraPrice,
		bid.Rnd,
		bid.TeamID.String(),
		bid.Signature,
		bid.State,
	)
	return err
}

func (b *Storage) UpdateBidState(auction uuid.UUID, extra_price int64, state BidState, stateExtra string) error {
	_, err := b.tx.Exec("UPDATE bids SET state=$1, state_extra=$2 WHERE auction=$3 AND extra_price=$4;", state, stateExtra, auction, extra_price)
	return err
}

func (b *Storage) UpdateBidPaymentID(auction uuid.UUID, extra_price int64, ID string) error {
	_, err := b.tx.Exec("UPDATE bids SET payment_id=$1 WHERE auction=$2 AND extra_price=$3;", ID, auction, extra_price)
	return err
}

func (b *Storage) UpdateBidPaymentUrl(auction uuid.UUID, extra_price int64, url string) error {
	_, err := b.tx.Exec("UPDATE bids SET payment_url=$1 WHERE auction=$2 AND extra_price=$3;", url, auction, extra_price)
	return err
}

func (b *Storage) UpdateBidPaymentDeadline(auction uuid.UUID, extra_price int64, deadline *big.Int) error {
	if deadline == nil {
		return errors.New("nil deadline")
	}
	_, err := b.tx.Exec("UPDATE bids SET payment_deadline=$1 WHERE auction=$2 AND extra_price=$3;", deadline.String(), auction, extra_price)
	return err
}

func (b *Storage) GetBidsOfAuction(auctionUUID uuid.UUID) ([]*Bid, error) {
	var bids []*Bid
	rows, err := b.tx.Query("SELECT extra_price, rnd, team_id, signature, state, state_extra, payment_id, payment_url, payment_deadline FROM bids WHERE auction=$1;", auctionUUID)
	if err != nil {
		return bids, err
	}
	defer rows.Close()
	for rows.Next() {
		var bid Bid
		var teamID sql.NullString
		var paymentDeadline sql.NullString
		err = rows.Scan(
			&bid.ExtraPrice,
			&bid.Rnd,
			&teamID,
			&bid.Signature,
			&bid.State,
			&bid.StateExtra,
			&bid.PaymentID,
			&bid.PaymentURL,
			&paymentDeadline,
		)
		if err != nil {
			return bids, err
		}
		bid.Auction = auctionUUID
		bid.TeamID, _ = new(big.Int).SetString(teamID.String, 10)
		bid.PaymentDeadline, _ = new(big.Int).SetString(paymentDeadline.String, 10)
		bids = append(bids, &bid)
	}
	return bids, nil
}
