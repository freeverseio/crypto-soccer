package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

type AuctionService struct {
	tx *sql.Tx
}

func NewAuctionService(tx *sql.Tx) *AuctionService {
	return &AuctionService{
		tx: tx,
	}
}

func (b AuctionService) PendingAuctions() ([]storage.Auction, error) {
	rows, err := b.tx.Query("SELECT id, player_id, currency_id, price, rnd, valid_until, signature, state, payment_url, state_extra, seller FROM auctions WHERE NOT (state = 'cancelled' OR state = 'failed' OR state = 'ended');")
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

func (b AuctionService) Auction(ID string) (*storage.Auction, error) {
	rows, err := b.tx.Query("SELECT player_id, currency_id, price, rnd, valid_until, signature, state, payment_url, state_extra, seller FROM auctions WHERE id = $1;", ID)
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

func (b AuctionService) Insert(auction storage.Auction) error {
	log.Debugf("[DBMS] + create Auction %v", b)
	_, err := b.tx.Exec("INSERT INTO auctions (id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, payment_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);",
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

func (b AuctionService) Update(auction storage.Auction) error {
	log.Debugf("[DBMS] + update Auction %v", b)
	_, err := b.tx.Exec(`UPDATE auctions SET 
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

func (b AuctionService) Bids(ID string) ([]storage.Bid, error) {
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

func (b AuctionService) BidInsert(bid storage.Bid) error {
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

func (b AuctionService) BidUpdate(bid storage.Bid) error {
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
