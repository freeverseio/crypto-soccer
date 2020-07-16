package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

type OfferService struct {
	tx *sql.Tx
}

func NewOfferService(tx *sql.Tx) *OfferService {
	return &OfferService{
		tx: tx,
	}
}

func (b OfferService) Bid() storage.BidService {
	return NewBidService(b.tx)
}

func (b OfferService) PendingOffers() ([]storage.Offer, error) {
	rows, err := b.tx.Query("SELECT id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, auction_id, team_id FROM offer WHERE NOT (state = 'cancelled' OR state = 'failed' OR state = 'ended') AND auction_id IS NOT NULL;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var offers []storage.Offer
	for rows.Next() {
		var offer storage.Offer
		err = rows.Scan(
			&offer.ID,
			&offer.PlayerID,
			&offer.CurrencyID,
			&offer.Price,
			&offer.Rnd,
			&offer.ValidUntil,
			&offer.Signature,
			&offer.State,
			&offer.StateExtra,
			&offer.Seller,
			&offer.Buyer,
			&offer.AuctionID,
			&offer.TeamID,
		)
		offers = append(offers, offer)
	}
	return offers, err
}

func (b OfferService) Offer(ID string) (*storage.Offer, error) {
	rows, err := b.tx.Query("SELECT player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, auction_id, team_id FROM auctions WHERE id = $1;", ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	var offer storage.Offer
	offer.ID = ID
	err = rows.Scan(
		&offer.PlayerID,
		&offer.CurrencyID,
		&offer.Price,
		&offer.Rnd,
		&offer.ValidUntil,
		&offer.Signature,
		&offer.State,
		&offer.StateExtra,
		&offer.Seller,
		&offer.Buyer,
		&offer.AuctionID,
		&offer.TeamID,
	)
	return &offer, err
}

func (b OfferService) Insert(offer storage.Offer) error {
	log.Debugf("[DBMS] + create Offer %v", b)
	_, err := b.tx.Exec("INSERT INTO offer (id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, team_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);",
		offer.ID,
		offer.PlayerID,
		offer.CurrencyID,
		offer.Price,
		offer.Rnd,
		offer.ValidUntil,
		offer.Signature,
		offer.State,
		offer.StateExtra,
		offer.Seller,
		offer.Buyer,
		offer.TeamID,
	)
	return err
}

func (b OfferService) Update(offer storage.Offer) error {
	log.Debugf("[DBMS] + update Offer %v", b)
	_, err := b.tx.Exec(`UPDATE offer SET 
		state=$1, 
		state_extra=$2,
		auction_id=$3,
		seller=$4
		WHERE id=$5;`,
		offer.State,
		offer.StateExtra,
		offer.AuctionID,
		offer.Seller,
		offer.ID,
	)
	return err
}
