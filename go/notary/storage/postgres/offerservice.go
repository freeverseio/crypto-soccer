package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (b StorageService) Offer(tx *sql.Tx, ID string) (*storage.Offer, error) {
	rows, err := tx.Query("SELECT player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, COALESCE(auction_id, ''), team_id FROM offers WHERE id = $1;", ID)
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

func (b StorageService) OfferByRndPrice(tx *sql.Tx, rnd int32, price int32) (*storage.Offer, error) {
	rows, err := tx.Query("SELECT id, player_id, currency_id, valid_until, signature, state, state_extra, seller, buyer, COALESCE(auction_id, ''), team_id FROM offers WHERE rnd = $1 AND price = $2;", rnd, price)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	var offer storage.Offer
	offer.Price = int64(price)
	offer.Rnd = int64(rnd)
	err = rows.Scan(
		&offer.ID,
		&offer.PlayerID,
		&offer.CurrencyID,
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

func (b StorageService) OfferByAuctionId(tx *sql.Tx, auctionId string) (*storage.Offer, error) {
	rows, err := tx.Query("SELECT player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, team_id FROM offers WHERE auction_id = $1;", auctionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	var offer storage.Offer
	offer.AuctionID = auctionId
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
		&offer.TeamID,
	)
	return &offer, err
}

func (b StorageService) OffersByPlayerId(tx *sql.Tx, playerId string) ([]storage.Offer, error) {
	rows, err := tx.Query("SELECT id, COALESCE(auction_id, ''), currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, team_id FROM offers WHERE player_id = $1;", playerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var offers []storage.Offer
	for rows.Next() {
		var offer storage.Offer
		offer.PlayerID = playerId
		err = rows.Scan(
			&offer.ID,
			&offer.AuctionID,
			&offer.CurrencyID,
			&offer.Price,
			&offer.Rnd,
			&offer.ValidUntil,
			&offer.Signature,
			&offer.State,
			&offer.StateExtra,
			&offer.Seller,
			&offer.Buyer,
			&offer.TeamID,
		)
		if err != nil {
			return offers, err
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

func (b StorageService) OfferInsert(tx *sql.Tx, offer storage.Offer) error {
	log.Debugf("[DBMS] + create Offer %v", b)
	_, err := tx.Exec("INSERT INTO offers (id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, team_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);",
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

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func (b StorageService) OfferUpdate(tx *sql.Tx, offer storage.Offer) error {
	log.Debugf("[DBMS] + update Offer %v", b)
	_, err := tx.Exec(`UPDATE offers SET 
		state=$1, 
		state_extra=$2,
		auction_id=$3,
		seller=$4
		WHERE id=$5;`,
		offer.State,
		offer.StateExtra,
		NewNullString(offer.AuctionID),
		NewNullString(offer.Seller),
		offer.ID,
	)
	return err
}
