package mockup

import (
	"database/sql"

	"github.com/freeverseio/crypto-soccer/go/notary/storage"
)

type StorageService struct {
	BeginFunc func() (storage.Tx, error)
}

type Tx struct {
	// Auction
	AuctionPendingAuctionsFunc func() ([]storage.Auction, error)
	AuctionFunc                func(ID string) (*storage.Auction, error)
	AuctionInsertFunc          func(auction storage.Auction) error
	AuctionUpdateFunc          func(auction storage.Auction) error
	AuctionsByPlayerIdFunc     func(ID string) ([]storage.Auction, error)

	// Bid
	BidFunc       func(auctionId string, extraPrice int64) (*storage.Bid, error)
	BidsFunc      func(auctionId string) ([]storage.Bid, error)
	BidInsertFunc func(bid storage.Bid) error
	BidUpdateFunc func(bid storage.Bid) error

	// PlayStore
	PlayStoreOrderFunc                   func(orderId string) (*storage.PlaystoreOrder, error)
	PlayStorePendingOrdersFunc           func() ([]storage.PlaystoreOrder, error)
	PlayStoreInsertFunc                  func(order storage.PlaystoreOrder) error
	PlayStoreUpdateStateFunc             func(order storage.PlaystoreOrder) error
	PlayStorePendingOrdersByPlayerIdFunc func(playerId string) ([]storage.PlaystoreOrder, error)

	// Offer
	OfferFunc              func(ID string) (*storage.Offer, error)
	OfferPendingOffersFunc func() ([]storage.Offer, error)
	OfferInsertFunc        func(offer storage.Offer) error
	OfferUpdateFunc        func(offer storage.Offer) error
	OfferByAuctionIdFunc   func(auctionId string) (*storage.Offer, error)
	OfferByRndPriceFunc    func(rnd int32, price int32) (*storage.Offer, error)
	OffersByPlayerIdFunc   func(playerId string) ([]storage.Offer, error)
}

func (b *StorageService) Begin() (storage.Tx, error) {
	return b.BeginFunc()
}
func (b *Tx) AuctionPendingAuctions(tx *sql.Tx) ([]storage.Auction, error) {
	return b.AuctionPendingAuctionsFunc()
}
func (b *Tx) Auction(ID string) (*storage.Auction, error) {
	return b.AuctionFunc(ID)
}
func (b *Tx) AuctionInsert(auction storage.Auction) error {
	return b.AuctionInsertFunc(auction)
}
func (b *Tx) AuctionUpdate(auction storage.Auction) error {
	return b.AuctionUpdateFunc(auction)
}
func (b *Tx) AuctionsByPlayerId(ID string) ([]storage.Auction, error) {
	return b.AuctionsByPlayerIdFunc(ID)
}
func (b *Tx) Bid(auctionId string, extraPrice int64) (*storage.Bid, error) {
	return b.BidFunc(auctionId, extraPrice)
}
func (b *Tx) Bids(auctionId string) ([]storage.Bid, error) {
	return b.BidsFunc(auctionId)
}
func (b *Tx) BidInsert(bid storage.Bid) error {
	return b.BidInsertFunc(bid)
}
func (b *Tx) BidUpdate(bid storage.Bid) error {
	return b.BidUpdateFunc(bid)
}
func (b *Tx) PlayStoreOrder(orderId string) (*storage.PlaystoreOrder, error) {
	return b.PlayStoreOrderFunc(orderId)
}
func (b *Tx) PlayStorePendingOrders(tx *sql.Tx) ([]storage.PlaystoreOrder, error) {
	return b.PlayStorePendingOrdersFunc()
}
func (b *Tx) PlayStoreInsert(order storage.PlaystoreOrder) error {
	return b.PlayStoreInsertFunc(order)
}
func (b *Tx) PlayStoreUpdateState(order storage.PlaystoreOrder) error {
	return b.PlayStoreUpdateStateFunc(order)
}
func (b *Tx) PlayStorePendingOrdersByPlayerId(playerId string) ([]storage.PlaystoreOrder, error) {
	return b.PlayStorePendingOrdersByPlayerIdFunc(playerId)
}
func (b *Tx) Offer(ID string) (*storage.Offer, error) {
	return b.OfferFunc(ID)
}
func (b *Tx) OfferPendingOffers() ([]storage.Offer, error) {
	return b.OfferPendingOffersFunc()
}
func (b *Tx) OfferInsert(offer storage.Offer) error {
	return b.OfferInsertFunc(offer)
}
func (b *Tx) OfferUpdate(offer storage.Offer) error {
	return b.OfferUpdateFunc(offer)
}
func (b *Tx) OfferByAuctionId(auctionId string) (*storage.Offer, error) {
	return b.OfferByAuctionIdFunc(auctionId)
}
func (b *Tx) OfferByRndPrice(rnd int32, price int32) (*storage.Offer, error) {
	return b.OfferByRndPriceFunc(rnd, price)
}
func (b *Tx) OffersByPlayerId(playerId string) ([]storage.Offer, error) {
	return b.OffersByPlayerIdFunc(playerId)
}
