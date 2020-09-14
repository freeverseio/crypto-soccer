package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

func (b *StorageHistoryTx) OfferInsert(offer storage.Offer) error {
	err := b.Tx.OfferInsert(offer)
	if err != nil {
		return err
	}
	return offerInsertHistory(b.Tx.tx, offer)
}

func (b *StorageHistoryTx) OfferUpdate(offer storage.Offer) error {
	currentOffer, err := b.Tx.Offer(offer.AuctionID)
	if err != nil {
		return err
	}
	if currentOffer == nil {
		return nil
	}
	if *currentOffer == offer {
		return nil
	}
	if err := b.Tx.OfferUpdate(offer); err != nil {
		return err
	}
	return offerInsertHistory(b.Tx.tx, offer)
}

func offerInsertHistory(tx *sql.Tx, offer storage.Offer) error {
	_, err := tx.Exec("INSERT INTO offers_histories_v2 (auction_id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, buyer_team_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);",
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
