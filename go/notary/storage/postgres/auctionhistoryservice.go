package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b StorageHistoryService) AuctionInsert(tx *sql.Tx, auction storage.Auction) error {
	if err := b.StorageService.AuctionInsert(tx, auction); err != nil {
		return err
	}
	return auctionInsertHistory(tx, auction)
}

func (b StorageHistoryService) AuctionUpdate(tx *sql.Tx, auction storage.Auction) error {
	if err := b.StorageService.AuctionUpdate(tx, auction); err != nil {
		return err
	}
	return auctionInsertHistory(tx, auction)
}

func auctionInsertHistory(tx *sql.Tx, auction storage.Auction) error {
	_, err := tx.Exec("INSERT INTO auctions_histories (id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, payment_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);",
		auction.ID,
		auction.PlayerID,
		auction.CurrencyID,
		auction.Price,
		auction.Rnd,
		auction.ValidUntil,
		auction.Signature,
		auction.State,
		auction.StateExtra,
		auction.Seller,
		auction.PaymentURL,
	)
	return err
}
