package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b StorageHistoryService) OfferInsert(tx *sql.Tx, offer storage.Offer) error {
	err := b.StorageService.OfferInsert(tx, offer)
	if err != nil {
		return err
	}
	return offerInsertHistory(tx, offer)
}

func (b StorageHistoryService) OfferUpdate(tx *sql.Tx, offer storage.Offer) error {
	if err := b.StorageService.OfferUpdate(tx, offer); err != nil {
		return err
	}
	return offerInsertHistory(tx, offer)
}

func offerInsertHistory(tx *sql.Tx, offer storage.Offer) error {
	_, err := tx.Exec("INSERT INTO offers_histories (id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, buyer_team_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);",
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
		offer.BuyerTeamID,
	)
	return err
}
