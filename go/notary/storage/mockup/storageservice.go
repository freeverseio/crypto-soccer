package mockup

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type StorageService struct {
	BeginFunc func() (*sql.Tx, error)
}

type Tx struct {
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
func (b *Tx) AuctionPendingAuctions(tx *sql.Tx) ([]storage.Auction, error) {
	return b.AuctionPendingAuctionsFunc(tx)
}
func (b *Tx) Auction(tx *sql.Tx, ID string) (*storage.Auction, error) {
	return b.AuctionFunc(tx, ID)
}
func (b *Tx) AuctionInsert(tx *sql.Tx, auction storage.Auction) error {
	return b.AuctionInsertFunc(tx, auction)
}
func (b *Tx) AuctionUpdate(tx *sql.Tx, auction storage.Auction) error {
	return b.AuctionUpdateFunc(tx, auction)
}
func (b *Tx) AuctionsByPlayerId(tx *sql.Tx, ID string) ([]storage.Auction, error) {
	return b.AuctionsByPlayerIdFunc(tx, ID)
}
func (b *Tx) Bid(tx *sql.Tx, auctionId string, extraPrice int64) (*storage.Bid, error) {
	return b.BidFunc(tx, auctionId, extraPrice)
}
func (b *Tx) Bids(tx *sql.Tx, auctionId string) ([]storage.Bid, error) {
	return b.BidsFunc(tx, auctionId)
}
func (b *Tx) BidInsert(tx *sql.Tx, bid storage.Bid) error {
	return b.BidInsertFunc(tx, bid)
}
func (b *Tx) BidUpdate(tx *sql.Tx, bid storage.Bid) error {
	return b.BidUpdateFunc(tx, bid)
}
func (b *Tx) PlayStoreOrder(tx *sql.Tx, orderId string) (*storage.PlaystoreOrder, error) {
	return b.PlayStoreOrderFunc(tx, orderId)
}
func (b *Tx) PlayStorePendingOrders(tx *sql.Tx) ([]storage.PlaystoreOrder, error) {
	return b.PlayStorePendingOrdersFunc(tx)
}
func (b *Tx) PlayStoreInsert(tx *sql.Tx, order storage.PlaystoreOrder) error {
	return b.PlayStoreInsertFunc(tx, order)
}
func (b *Tx) PlayStoreUpdateState(tx *sql.Tx, order storage.PlaystoreOrder) error {
	return b.PlayStoreUpdateStateFunc(tx, order)
}
func (b *Tx) PlayStorePendingOrdersByPlayerId(tx *sql.Tx, playerId string) ([]storage.PlaystoreOrder, error) {
	return b.PlayStorePendingOrdersByPlayerIdFunc(tx, playerId)
}
func (b *Tx) Offer(tx *sql.Tx, ID string) (*storage.Offer, error) {
	return b.OfferFunc(tx, ID)
}
func (b *Tx) OfferPendingOffers(tx *sql.Tx) ([]storage.Offer, error) {
	return b.OfferPendingOffersFunc(tx)
}
func (b *Tx) OfferInsert(tx *sql.Tx, offer storage.Offer) error {
	return b.OfferInsertFunc(tx, offer)
}
func (b *Tx) OfferUpdate(tx *sql.Tx, offer storage.Offer) error {
	return b.OfferUpdateFunc(tx, offer)
}
func (b *Tx) OfferByAuctionId(tx *sql.Tx, auctionId string) (*storage.Offer, error) {
	return b.OfferByAuctionIdFunc(tx, auctionId)
}
func (b *Tx) OfferByRndPrice(tx *sql.Tx, rnd int32, price int32) (*storage.Offer, error) {
	return b.OfferByRndPriceFunc(tx, rnd, price)
}
func (b *Tx) OffersByPlayerId(tx *sql.Tx, playerId string) ([]storage.Offer, error) {
	return b.OffersByPlayerIdFunc(tx, playerId)
}
