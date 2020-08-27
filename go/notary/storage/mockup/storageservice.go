package mockup

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type StorageService struct {
	BeginFunc func() (*sql.Tx, error)

	// Auction
	AuctionPendingAuctionsFunc func(tx *sql.Tx) ([]storage.Auction, error)
	AuctionFunc                func(tx *sql.Tx, ID string) (*storage.Auction, error)
	AuctionInsertFunc          func(tx *sql.Tx, auction storage.Auction) error
	AuctionUpdateFunc          func(tx *sql.Tx, auction storage.Auction) error
	AuctionsByPlayerIdFunc     func(tx *sql.Tx, ID string) ([]storage.Auction, error)

	// Bid
	BidFunc       func(tx *sql.Tx, auctionId string, extraPrice int64) (*storage.Bid, error)
	BidsFunc      func(tx *sql.Tx, auctionId string) ([]storage.Bid, error)
	BidInsertFunc func(tx *sql.Tx, bid storage.Bid) error
	BidUpdateFunc func(tx *sql.Tx, bid storage.Bid) error

	// PlayStore
	PlayStoreOrderFunc                   func(tx *sql.Tx, orderId string) (*storage.PlaystoreOrder, error)
	PlayStorePendingOrdersFunc           func(tx *sql.Tx) ([]storage.PlaystoreOrder, error)
	PlayStoreInsertFunc                  func(tx *sql.Tx, order storage.PlaystoreOrder) error
	PlayStoreUpdateStateFunc             func(tx *sql.Tx, order storage.PlaystoreOrder) error
	PlayStorePendingOrdersByPlayerIdFunc func(tx *sql.Tx, playerId string) ([]storage.PlaystoreOrder, error)

	// Offer
	OfferFunc              func(tx *sql.Tx, ID string) (*storage.Offer, error)
	OfferPendingOffersFunc func(tx *sql.Tx) ([]storage.Offer, error)
	OfferInsertFunc        func(tx *sql.Tx, offer storage.Offer) error
	OfferUpdateFunc        func(tx *sql.Tx, offer storage.Offer) error
	OfferByAuctionIdFunc   func(tx *sql.Tx, auctionId string) (*storage.Offer, error)
	OfferByRndPriceFunc    func(tx *sql.Tx, rnd int32, price int32) (*storage.Offer, error)
	OffersByPlayerIdFunc   func(tx *sql.Tx, playerId string) ([]storage.Offer, error)
}

func (b *StorageService) Begin() (*sql.Tx, error) {
	return b.BeginFunc()
}
func (b *StorageService) AuctionPendingAuctions(tx *sql.Tx) ([]storage.Auction, error) {
	return b.AuctionPendingAuctionsFunc(tx)
}
func (b *StorageService) Auction(tx *sql.Tx, ID string) (*storage.Auction, error) {
	return b.AuctionFunc(tx, ID)
}
func (b *StorageService) AuctionInsert(tx *sql.Tx, auction storage.Auction) error {
	return b.AuctionInsertFunc(tx, auction)
}
func (b *StorageService) AuctionUpdate(tx *sql.Tx, auction storage.Auction) error {
	return b.AuctionUpdateFunc(tx, auction)
}
func (b *StorageService) AuctionsByPlayerId(tx *sql.Tx, ID string) ([]storage.Auction, error) {
	return b.AuctionsByPlayerIdFunc(tx, ID)
}
func (b *StorageService) Bid(tx *sql.Tx, auctionId string, extraPrice int64) (*storage.Bid, error) {
	return b.BidFunc(tx, auctionId, extraPrice)
}
func (b *StorageService) Bids(tx *sql.Tx, auctionId string) ([]storage.Bid, error) {
	return b.BidsFunc(tx, auctionId)
}
func (b *StorageService) BidInsert(tx *sql.Tx, bid storage.Bid) error {
	return b.BidInsertFunc(tx, bid)
}
func (b *StorageService) BidUpdate(tx *sql.Tx, bid storage.Bid) error {
	return b.BidUpdateFunc(tx, bid)
}
func (b *StorageService) PlayStoreOrder(tx *sql.Tx, orderId string) (*storage.PlaystoreOrder, error) {
	return b.PlayStoreOrderFunc(tx, orderId)
}
func (b *StorageService) PlayStorePendingOrders(tx *sql.Tx) ([]storage.PlaystoreOrder, error) {
	return b.PlayStorePendingOrdersFunc(tx)
}
func (b *StorageService) PlayStoreInsert(tx *sql.Tx, order storage.PlaystoreOrder) error {
	return b.PlayStoreInsertFunc(tx, order)
}
func (b *StorageService) PlayStoreUpdateState(tx *sql.Tx, order storage.PlaystoreOrder) error {
	return b.PlayStoreUpdateStateFunc(tx, order)
}
func (b *StorageService) PlayStorePendingOrdersByPlayerId(tx *sql.Tx, playerId string) ([]PlaystoreOrder, error) {
	return b.PlayStorePendingOrdersByPlayerIdFunc(tx, playerId)
}
func (b *StorageService) Offer(tx *sql.Tx, ID string) (*storage.Offer, error) {
	return b.OfferFunc(tx, ID)
}
func (b *StorageService) OfferPendingOffers(tx *sql.Tx) ([]storage.Offer, error) {
	return b.OfferPendingOffersFunc(tx)
}
func (b *StorageService) OfferInsert(tx *sql.Tx, offer storage.Offer) error {
	return b.OfferInsertFunc(tx, offer)
}
func (b *StorageService) OfferUpdate(tx *sql.Tx, offer storage.Offer) error {
	return b.OfferUpdateFunc(tx, offer)
}
func (b *StorageService) OfferByAuctionId(tx *sql.Tx, auctionId string) (*storage.Offer, error) {
	return b.OfferByAuctionIdFunc(tx, auctionId)
}
func (b *StorageService) OfferByRndPrice(tx *sql.Tx, rnd int32, price int32) (*storage.Offer, error) {
	return b.OfferByRndPriceFunc(tx, rnd, price)
}
func (b *StorageService) OffersByPlayerId(tx *sql.Tx, playerId string) ([]storage.Offer, error) {
	return b.OffersByPlayerIdFunc(tx, playerId)
}
