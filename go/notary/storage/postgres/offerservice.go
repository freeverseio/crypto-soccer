package postgres

import (
	"database/sql"
	"errors"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
	log "github.com/sirupsen/logrus"
)

func (b *Tx) OfferPendingOffers() ([]storage.Offer, error) {
	rows, err := b.tx.Query("SELECT auction_id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, buyer_team_id FROM offers_v2 WHERE state = 'started';")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var offers []storage.Offer
	for rows.Next() {
		var offer storage.Offer
		err = rows.Scan(
			&offer.AuctionID,
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
			&offer.BuyerTeamID,
		)
		offers = append(offers, offer)
	}
	return offers, err
}

func (b *Tx) Offer(AuctionID string) (*storage.Offer, error) {
	rows, err := b.tx.Query("SELECT player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, buyer_team_id FROM offers_v2 WHERE auction_id = $1;", AuctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.New("Could not find the offer you queried by auctionID")
	}
	var offer storage.Offer
	offer.AuctionID = AuctionID
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
		&offer.BuyerTeamID,
	)
	return &offer, err
}

func (b *Tx) OfferByRndPrice(rnd int32, price int32) (*storage.Offer, error) {
	rows, err := b.tx.Query("SELECT auction_id, player_id, currency_id, valid_until, signature, state, state_extra, seller, buyer), buyer_team_id FROM offers_v2 WHERE rnd = $1 AND price = $2;", rnd, price)
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
		&offer.AuctionID,
		&offer.PlayerID,
		&offer.CurrencyID,
		&offer.ValidUntil,
		&offer.Signature,
		&offer.State,
		&offer.StateExtra,
		&offer.Seller,
		&offer.Buyer,
		&offer.BuyerTeamID,
	)
	return &offer, err
}

func (b *Tx) OfferByAuctionId(auctionId string) (*storage.Offer, error) {
	rows, err := b.tx.Query("SELECT auction_id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, buyer_team_id FROM offers_v2 WHERE auction_id = $1;", auctionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}
	var offer storage.Offer
	err = rows.Scan(
		&offer.AuctionID,
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
		&offer.BuyerTeamID,
	)
	return &offer, err
}

func (b *Tx) OffersByPlayerId(playerId string) ([]storage.Offer, error) {
	rows, err := b.tx.Query("SELECT auction_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, buyer_team_id FROM offers_v2 WHERE player_id = $1;", playerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var offers []storage.Offer
	for rows.Next() {
		var offer storage.Offer
		offer.PlayerID = playerId
		err = rows.Scan(
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
			&offer.BuyerTeamID,
		)
		if err != nil {
			return offers, err
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

func (b *Tx) OfferInsert(offer storage.Offer) error {
	log.Debugf("[DBMS] + create Offer %v", b)
	if offer.AuctionID == "" {
		return errors.New("Trying to insert an auction with empty auctionID")
	}
	_, err := b.tx.Exec("INSERT INTO offers_v2 (auction_id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, buyer_team_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);",
		offer.AuctionID,
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
		offer.BuyerTeamID,
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

func (b *Tx) OfferUpdate(offer storage.Offer) error {
	log.Debugf("[DBMS] + update Offer %v", b)
	log.Warning(offer.State)
	log.Warning(offer.StateExtra)
	log.Warning(offer.Seller)
	log.Warning(offer.AuctionID)

	rows, err := b.tx.Exec(`UPDATE offers_v2 SET 
		state=$1, 
		state_extra=$2,
		seller=$3
		WHERE auction_id=$4;`,
		offer.State,
		offer.StateExtra,
		NewNullString(offer.Seller),
		offer.AuctionID,
	)
	if err != nil {
		return err
	}
	nInserted, err := rows.RowsAffected()
	if nInserted == 0 {
		return errors.New("UPDATE: could not find an offer to update")
	}
	return err
}

func (b *Tx) OfferCancel(AuctionID string) error {
	log.Debugf("[DBMS] + update Offer %v", b)
	_, err := b.tx.Exec(`UPDATE offers_v2 SET 
		state=$1, 
		WHERE auction_id=$2;`,
		storage.OfferCancelled,
		AuctionID,
	)
	return err
}
