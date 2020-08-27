package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *StorageHistoryService) BidInsert(bid storage.Bid) error {
	if err := b.StorageService.BidInsert(bid); err != nil {
		return err
	}
	return bidInsertHistory(b.StorageService.tx, bid)
}

func (b *StorageHistoryService) BidUpdate(bid storage.Bid) error {
	currentBids, err := b.StorageService.Bid(bid.AuctionID, bid.ExtraPrice)
	if err != nil {
		return err
	}
	if currentBids == nil {
		return nil
	}
	if *currentBids == bid {
		return nil
	}
	if err := b.StorageService.BidUpdate(bid); err != nil {
		return err
	}
	return bidInsertHistory(b.StorageService.tx, bid)
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
