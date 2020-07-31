package postgres

// import (
// 	"database/sql"

// 	"github.com/freeverseio/crypto-soccer/go/notary/storage"
// )

// type BidHistoryService struct {
// 	BidService
// }

// func NewBidHistoryService(tx *sql.Tx) *BidHistoryService {
// 	return &BidHistoryService{*NewBidService(tx)}
// }

// func (b BidHistoryService) Insert(bid storage.Bid) error {
// 	if err := b.BidService.Insert(bid); err != nil {
// 		return err
// 	}
// 	return b.insertHistory(bid)
// }

// func (b BidHistoryService) Update(bid storage.Bid) error {
// 	if err := b.BidService.Update(bid); err != nil {
// 		return err
// 	}
// 	return b.insertHistory(bid)
// }

// func (b BidHistoryService) insertHistory(bid storage.Bid) error {
// 	_, err := b.tx.Exec(`INSERT INTO bids_histories
// 			(auction_id,
// 			extra_price,
// 			rnd,
// 			team_id,
// 			signature,
// 			state,
// 			state_extra,
// 			payment_id,
// 			payment_url,
// 			payment_deadline)
// 			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
// 		bid.AuctionID,
// 		bid.ExtraPrice,
// 		bid.Rnd,
// 		bid.TeamID,
// 		bid.Signature,
// 		bid.State,
// 		bid.StateExtra,
// 		bid.PaymentID,
// 		bid.PaymentURL,
// 		bid.PaymentDeadline,
// 	)
// 	return err
// }
