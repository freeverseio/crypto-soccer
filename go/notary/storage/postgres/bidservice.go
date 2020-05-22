package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

type BidService struct {
	tx *sql.Tx
}

func NewBidService(tx *sql.Tx) storage.BidService {
	return &BidService{
		tx: tx,
	}
}

func (b BidService) Bids(ID string) ([]storage.Bid, error) {
	rows, err := b.tx.Query("SELECT extra_price, rnd, team_id, signature, state, state_extra, payment_id, payment_url, payment_deadline FROM bids WHERE auction_id=$1;", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []storage.Bid
	for rows.Next() {
		var bid storage.Bid
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

func (b BidService) Insert(bid storage.Bid) error {
	log.Debugf("[DBMS] + create Bid %v", b)
	_, err := b.tx.Exec(`INSERT INTO bids 
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
		bid.AuctionID,
		bid.ExtraPrice,
		bid.Rnd,
		bid.TeamID,
		bid.Signature,
		bid.State,
		bid.StateExtra,
		bid.PaymentID,
		bid.PaymentURL,
		bid.PaymentDeadline,
	)
	return err
}

func (b BidService) Update(bid storage.Bid) error {
	log.Debugf("[DBMS] + update Bid %v", b)
	_, err := b.tx.Exec(`UPDATE bids SET 
		state=$1, 
		state_extra=$2,
		payment_id=$3,
		payment_url=$4,
		payment_deadline=$5
		WHERE auction_id=$6 AND extra_price=$7;`,
		bid.State,
		bid.StateExtra,
		bid.PaymentID,
		bid.PaymentURL,
		bid.PaymentDeadline,
		bid.AuctionID,
		bid.ExtraPrice,
	)
	return err
}
