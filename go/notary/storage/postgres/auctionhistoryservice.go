package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *StorageHistoryTx) AuctionInsert(auction storage.Auction) error {
	if err := b.Tx.AuctionInsert(auction); err != nil {
		return err
	}
	return auctionInsertHistory(b.Tx.tx, auction)
}

func (b *StorageHistoryTx) AuctionUpdate(auction storage.Auction) error {
	currentAuction, err := b.Tx.Auction(auction.ID)
	if err != nil {
		return err
	}
	if currentAuction == nil {
		return nil
	}
	if *currentAuction == auction {
		return nil
	}
	if err := b.Tx.AuctionUpdate(auction); err != nil {
		return err
	}
	return auctionInsertHistory(b.Tx.tx, auction)
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
