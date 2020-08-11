package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (b StorageService) AuctionPendingAuctions(tx *sql.Tx) ([]storage.Auction, error) {
	rows, err := tx.Query("SELECT id, player_id, currency_id, price, rnd, valid_until, signature, state, payment_url, state_extra, seller FROM auctions WHERE NOT (state = 'cancelled' OR state = 'failed' OR state = 'ended');")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var auctions []storage.Auction
	for rows.Next() {
		var auction storage.Auction
		err = rows.Scan(
			&auction.ID,
			&auction.PlayerID,
			&auction.CurrencyID,
			&auction.Price,
			&auction.Rnd,
			&auction.ValidUntil,
			&auction.Signature,
			&auction.State,
			&auction.PaymentURL,
			&auction.StateExtra,
			&auction.Seller,
		)
		auctions = append(auctions, auction)
	}
	return auctions, err
}

func (b StorageService) Auction(tx *sql.Tx, ID string) (*storage.Auction, error) {
	rows, err := tx.Query("SELECT player_id, currency_id, price, rnd, valid_until, signature, state, payment_url, state_extra, seller FROM auctions WHERE id = $1;", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	var auction storage.Auction
	auction.ID = ID
	err = rows.Scan(
		&auction.PlayerID,
		&auction.CurrencyID,
		&auction.Price,
		&auction.Rnd,
		&auction.ValidUntil,
		&auction.Signature,
		&auction.State,
		&auction.PaymentURL,
		&auction.StateExtra,
		&auction.Seller,
	)
	return &auction, err
}

func (b StorageService) AuctionsByPlayerId(tx *sql.Tx, ID string) ([]storage.Auction, error) {
	rows, err := tx.Query("SELECT id, currency_id, price, rnd, valid_until, signature, state, payment_url, state_extra, seller FROM auctions WHERE player_id = $1;", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var auctions []storage.Auction
	for rows.Next() {
		var auction storage.Auction
		auction.PlayerID = ID
		err = rows.Scan(
			&auction.ID,
			&auction.CurrencyID,
			&auction.Price,
			&auction.Rnd,
			&auction.ValidUntil,
			&auction.Signature,
			&auction.State,
			&auction.PaymentURL,
			&auction.StateExtra,
			&auction.Seller,
		)
		if err != nil {
			return auctions, err
		}
		auctions = append(auctions, auction)
	}

	return auctions, err
}

func (b StorageService) AuctionInsert(tx *sql.Tx, auction storage.Auction) error {
	log.Debugf("[DBMS] + create Auction %v", b)
	_, err := tx.Exec("INSERT INTO auctions (id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, payment_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);",
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

func (b StorageService) AuctionUpdate(tx *sql.Tx, auction storage.Auction) error {
	log.Debugf("[DBMS] + update Auction %v", b)
	_, err := tx.Exec(`UPDATE auctions SET 
		state=$1, 
		state_extra=$2,
		payment_url=$3
		WHERE id=$4;`,
		auction.State,
		auction.StateExtra,
		auction.PaymentURL,
		auction.ID,
	)
	return err
}
