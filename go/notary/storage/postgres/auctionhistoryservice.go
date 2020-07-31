package postgres

// type AuctionHistoryService struct {
// 	AuctionService
// }

// func NewAuctionHistoryService(tx *sql.Tx) *AuctionHistoryService {
// 	return &AuctionHistoryService{*NewAuctionService(tx)}
// }

// func (b AuctionHistoryService) Bid() storage.BidService {
// 	return *NewBidHistoryService(b.tx)
// }

// func (b AuctionHistoryService) Insert(auction storage.Auction) error {
// 	if err := b.AuctionService.Insert(auction); err != nil {
// 		return err
// 	}
// 	return b.insertHistory(auction)
// }

// func (b AuctionHistoryService) Update(auction storage.Auction) error {
// 	if err := b.AuctionService.Update(auction); err != nil {
// 		return err
// 	}
// 	return b.insertHistory(auction)
// }

// func (b AuctionHistoryService) insertHistory(auction storage.Auction) error {
// 	_, err := b.tx.Exec("INSERT INTO auctions_histories (id, player_id, currency_id, price, rnd, valid_until, signature, state, state_extra, seller, payment_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);",
// 		auction.ID,
// 		auction.PlayerID,
// 		auction.CurrencyID,
// 		auction.Price,
// 		auction.Rnd,
// 		auction.ValidUntil,
// 		auction.Signature,
// 		auction.State,
// 		auction.StateExtra,
// 		auction.Seller,
// 		auction.PaymentURL,
// 	)
// 	return err
// }
