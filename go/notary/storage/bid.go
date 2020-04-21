package storage

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type BidState string

const (
	BIDREFUSED  BidState = "REFUSED"
	BidAccepted BidState = "accepted"
	BidPaid     BidState = "paid"
	BidPaying   BidState = "paying"
	BidFailed   BidState = "failed"
)

type Bid struct {
	AuctionID       string
	ExtraPrice      int64
	Rnd             int64
	TeamID          string
	Signature       string
	State           BidState
	StateExtra      string
	PaymentID       string
	PaymentURL      string
	PaymentDeadline int64
}

func NewBid() *Bid {
	bid := Bid{}
	bid.State = BidAccepted
	return &bid
}

func BidsByAuctionID(tx *sql.Tx, ID string) ([]Bid, error) {
	rows, err := tx.Query("SELECT extra_price, rnd, team_id, signature, state, state_extra, payment_id, payment_url, payment_deadline FROM bids WHERE auction_id=$1;", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []Bid
	for rows.Next() {
		var bid Bid
		bid.AuctionID = ID
		err = rows.Scan(
			&bid.ExtraPrice,
			&bid.Rnd,
			&bid.TeamID,
			&bid.Signature,
			&bid.State,
			&bid.StateExtra,
			&bid.PaymentID,
			&bid.PaymentURL,
			&bid.PaymentDeadline,
		)
		if err != nil {
			return bids, err
		}
		bids = append(bids, bid)
	}
	return bids, nil
}

func (b Bid) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] + create Bid %v", b)
	_, err := tx.Exec(`INSERT INTO bids 
			(auction_id, 
			extra_price,
			rnd, team_id, 
			signature, 
			state,
			state_extra,
			payment_id,
			payment_url,
			payment_deadline) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		b.AuctionID,
		b.ExtraPrice,
		b.Rnd,
		b.TeamID,
		b.Signature,
		b.State,
		b.StateExtra,
		b.PaymentID,
		b.PaymentURL,
		b.PaymentDeadline,
	)
	return err
}

func (b Bid) Update(tx *sql.Tx) error {
	log.Debugf("[DBMS] + update Bid %v", b)
	_, err := tx.Exec(`UPDATE bids SET 
		state=$1, 
		state_extra=$2,
		payment_id=$3,
		payment_url=$4,
		payment_deadline=$5
		WHERE auction_id=$6 AND extra_price=$7;`,
		b.State,
		b.StateExtra,
		b.PaymentID,
		b.PaymentURL,
		b.PaymentDeadline,
		b.AuctionID,
		b.ExtraPrice,
	)
	return err
}

// func (b *Storage) CreateBid(bid Bid) error {
// 	log.Infof("[DBMS] + create Bid %v", bid)
// 	_, err := b.db.Exec("INSERT INTO bids (auction, extra_price, rnd, team_id, signature, state) VALUES ($1, $2, $3, $4, $5, $6);",
// 		bid.Auction,
// 		bid.ExtraPrice,
// 		bid.Rnd,
// 		bid.TeamID.String(),
// 		bid.Signature,
// 		bid.State,
// 	)
// 	return err
// }

// func (b *Storage) UpdateBidState(auction uuid.UUID, extra_price int64, state BidState, stateExtra string) error {
// 	_, err := b.db.Exec("UPDATE bids SET state=$1, state_extra=$2 WHERE auction=$3 AND extra_price=$4;", state, stateExtra, auction, extra_price)
// 	return err
// }

// func (b *Storage) UpdateBidPaymentID(auction uuid.UUID, extra_price int64, ID string) error {
// 	_, err := b.db.Exec("UPDATE bids SET payment_id=$1 WHERE auction=$2 AND extra_price=$3;", ID, auction, extra_price)
// 	return err
// }

// func (b *Storage) UpdateBidPaymentUrl(auction uuid.UUID, extra_price int64, url string) error {
// 	_, err := b.db.Exec("UPDATE bids SET payment_url=$1 WHERE auction=$2 AND extra_price=$3;", url, auction, extra_price)
// 	return err
// }

// func (b *Storage) UpdateBidPaymentDeadline(auction uuid.UUID, extra_price int64, deadline *big.Int) error {
// 	if deadline == nil {
// 		return errors.New("nil deadline")
// 	}
// 	_, err := b.db.Exec("UPDATE bids SET payment_deadline=$1 WHERE auction=$2 AND extra_price=$3;", deadline.String(), auction, extra_price)
// 	return err
// }

// func (b *Storage) GetBidsOfAuction(auctionUUID uuid.UUID) ([]*Bid, error) {
// 	var bids []*Bid
// 	rows, err := b.db.Query("SELECT extra_price, rnd, team_id, signature, state, state_extra, payment_id, payment_url, payment_deadline FROM bids WHERE auction=$1;", auctionUUID)
// 	if err != nil {
// 		return bids, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var bid Bid
// 		var teamID sql.NullString
// 		err = rows.Scan(
// 			&bid.ExtraPrice,
// 			&bid.Rnd,
// 			&teamID,
// 			&bid.Signature,
// 			&bid.State,
// 			&bid.StateExtra,
// 			&bid.PaymentID,
// 			&bid.PaymentURL,
// 			&bid.PaymentDeadline,
// 		)
// 		if err != nil {
// 			return bids, err
// 		}
// 		bid.Auction = auctionUUID
// 		bid.TeamID, _ = new(big.Int).SetString(teamID.String, 10)
// 		bids = append(bids, &bid)
// 	}
// 	return bids, nil
// }
