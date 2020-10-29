package storage

type StorageService interface {
	Begin() (Tx, error)
}

type Tx interface {
	Rollback() error
	Commit() error

	// Auction
	AuctionPendingAuctions() ([]Auction, error)
	AuctionPendingOrderlessAuctions() ([]Auction, error)
	Auction(ID string) (*Auction, error)
	AuctionInsert(auction Auction) error
	AuctionUpdate(auction Auction) error
	AuctionCancel(ID string) error
	AuctionsByPlayerId(ID string) ([]Auction, error)

	// Bid
	Bid(auctionId string, extraPrice int64) (*Bid, error)
	Bids(auctionId string) ([]Bid, error)
	BidInsert(bid Bid) error
	BidUpdate(bid Bid) error

	// PlayStore
	PlayStoreOrder(orderId string) (*PlaystoreOrder, error)
	PlayStorePendingOrders() ([]PlaystoreOrder, error)
	PlayStoreInsert(order PlaystoreOrder) error
	PlayStoreUpdateState(order PlaystoreOrder) error
	PlayStorePendingOrdersByPlayerId(playerId string) ([]PlaystoreOrder, error)

	// Offer
	Offer(ID string) (*Offer, error)
	OfferPendingOffers() ([]Offer, error)
	OfferInsert(offer Offer) error
	OfferUpdate(offer Offer) error
	OfferCancel(ID string) error
	OfferByRndPrice(rnd int32, price int32) (*Offer, error)
	OffersByPlayerId(playerId string) ([]Offer, error)
	OffersStartedByPlayerId(playerId string) ([]Offer, error)
	CancelAllOffersByPlayerId(playerId string) error
}
