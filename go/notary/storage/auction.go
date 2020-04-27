package storage

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type AuctionState string

const (
	AuctionAssetFrozen        AuctionState = "asset_frozen"
	AuctionPaying             AuctionState = "paying"
	AuctionWithdrableBySeller AuctionState = "withadrable_by_seller"
	AuctionWithdrableByBuyer  AuctionState = "withadrable_by_buyer"
	AuctionStarted            AuctionState = "started"
	AuctionEnded              AuctionState = "ended"
	AuctionCancelled          AuctionState = "cancelled"
	AuctionFailed             AuctionState = "failed"
)

type Auction struct {
	ID         string
	PlayerID   string
	CurrencyID int
	Price      int64
	Rnd        int64
	ValidUntil int64
	Signature  string
	State      AuctionState
	StateExtra string
	PaymentURL string
	Seller     string
}

func NewAuction() *Auction {
	auction := Auction{}
	auction.State = AuctionStarted
	return &auction
}

func PendingAuctions(tx *sql.Tx) ([]Auction, error) {
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
