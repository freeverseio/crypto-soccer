package postgres

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type OfferHistoryService struct {
	OfferService
}

func NewOfferHistoryService(tx *sql.Tx) *OfferHistoryService {
	return &OfferHistoryService{*NewOfferService(tx)}
}

// func (b OfferHistoryService) Bid() storage.BidService {
// 	return *NewBidHistoryService(b.tx)
// }

func (b OfferHistoryService) Insert(offer storage.Offer) (int64, error) {
	id, err := b.OfferService.Insert(offer)
	if err != nil {
		return 0, err
	}
	offer.ID = id
	return id, b.insertHistory(offer)
}

func (b OfferHistoryService) Update(offer storage.Offer) error {
	if err := b.OfferService.Update(offer); err != nil {
		return err
	}
	return b.insertHistory(offer)
}

func (b OfferHistoryService) insertHistory(offer storage.Offer) error {
	_, err := b.tx.Exec("INSERT INTO offers_histories (id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, buyer, team_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);",
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
