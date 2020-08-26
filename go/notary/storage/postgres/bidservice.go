package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (b StorageService) Bid(tx *sql.Tx, auctionId string, extraPrice int64) (*storage.Bid, error) {
	rows, err := tx.Query("SELECT rnd, team_id, signature, state, state_extra, payment_id, payment_url, payment_deadline FROM bids WHERE auction_id=$1 AND extra_price=$2;", auctionId, extraPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var bid storage.Bid
	bid.AuctionID = auctionId
	bid.ExtraPrice = extraPrice
	err = rows.Scan(
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
		return nil, err
	}
	return &bid, nil
}

func (b StorageService) Bids(tx *sql.Tx, ID string) ([]storage.Bid, error) {
	rows, err := tx.Query("SELECT extra_price, rnd, team_id, signature, state, state_extra, payment_id, payment_url, payment_deadline FROM bids WHERE auction_id=$1;", ID)
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

func (b StorageService) BidInsert(tx *sql.Tx, bid storage.Bid) error {
	log.Debugf("[DBMS] + create Bid %v", b)
	_, err := tx.Exec(`INSERT INTO bids 
			(auction_id, 
			extra_price,
			rnd, 
			team_id, 
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

func (b StorageService) BidUpdate(tx *sql.Tx, bid storage.Bid) error {
	log.Debugf("[DBMS] + update Bid %v", b)
	_, err := tx.Exec(`UPDATE bids SET 
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
