package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	"google.golang.org/appengine/log"
)

type AuctionService struct {
	tx *sql.Tx
}

func NewAuctionService(tx *sql.Tx) *AuctionService {
	return &AuctionService{
		tx: tx,
	}
}

func PendingAuctions(tx *sql.Tx) ([]storage.Auction, error) {
	rows, err := tx.Query("SELECT id, player_id, currency_id, price, rnd, valid_until, signature, state, payment_url, state_extra, seller FROM auctions WHERE NOT (state = 'cancelled' OR state = 'failed' OR state = 'ended');")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var auctions []Auction
	for rows.Next() {
		var auction Auction
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

func AuctionByID(tx *sql.Tx, ID string) (*Auction, error) {
	rows, err := tx.Query("SELECT player_id, currency_id, price, rnd, valid_until, signature, state, payment_url, state_extra, seller FROM auctions WHERE id = $1;", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	var auction Auction
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

func (b Auction) Insert(tx *sql.Tx) error {
	log.Debugf("[DBMS] + create Auction %v", b)
	_, err := tx.Exec("INSERT INTO auctions (id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, payment_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);",
		b.ID,
		b.PlayerID,
		b.CurrencyID,
		b.Price,
		b.Rnd,
		b.ValidUntil,
		b.Signature,
		b.State,
		b.StateExtra,
		b.Seller,
		b.PaymentURL,
	)
	return err
}

func (b Auction) Update(tx *sql.Tx) error {
	log.Debugf("[DBMS] + update Auction %v", b)
	_, err := tx.Exec(`UPDATE auctions SET 
		state=$1, 
		state_extra=$2,
		payment_url=$3
		WHERE id=$4;`,
		b.State,
		b.StateExtra,
		b.PaymentURL,
		b.ID,
	)
	return err
}
