package storage

import "database/sql"

type StorageService interface {
	DB() *sql.DB

	// Auction
	AuctionPendingAuctions(tx *sql.Tx) ([]Auction, error)
	Auction(tx *sql.Tx, ID string) (*Auction, error)
	AuctionInsert(tx *sql.Tx, auction Auction) error
	AuctionUpdate(tx *sql.Tx, auction Auction) error

	// Bid
	Bids(tx *sql.Tx, auctionId string) ([]Bid, error)
	BidInsert(tx *sql.Tx, bid Bid) error
	BidUpdate(tx *sql.Tx, bid Bid) error

	// PlayStore
	PlayStoreOrder(tx *sql.Tx, orderId string) (*PlaystoreOrder, error)
	PlayStorePendingOrders(tx *sql.Tx) ([]PlaystoreOrder, error)
	PlayStoreInsert(tx *sql.Tx, order PlaystoreOrder) error
	PlayStoreUpdateState(tx *sql.Tx, order PlaystoreOrder) error
	PlayStorePendingOrdersByPlayerId(tx *sql.Tx, playerId string) ([]PlaystoreOrder, error)

	// Offer
	Offer(tx *sql.Tx, ID int64) (*Offer, error)
	OfferInsert(tx *sql.Tx, offer Offer) (int64, error)
	OfferUpdate(tx *sql.Tx, offer Offer) error
	OfferByAuctionId(tx *sql.Tx, auctionId string) (*Offer, error)
	OfferByRndPrice(tx *sql.Tx, rnd int32, price int32) (*Offer, error)
}
