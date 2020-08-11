package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b StorageHistoryService) BidInsert(tx *sql.Tx, bid storage.Bid) error {
	if err := b.StorageService.BidInsert(tx, bid); err != nil {
		return err
	}
	return bidInsertHistory(tx, bid)
}

func (b StorageHistoryService) BidUpdate(tx *sql.Tx, bid storage.Bid) error {
	if err := b.StorageService.BidUpdate(tx, bid); err != nil {
		return err
	}
	return bidInsertHistory(tx, bid)
}

func bidInsertHistory(tx *sql.Tx, bid storage.Bid) error {
	_, err := tx.Exec(`INSERT INTO bids_histories
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
