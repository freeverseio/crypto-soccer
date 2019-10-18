package storage

type AuctionState string

const (
	AUCTION_STARTED      AuctionState = "STARTED"
	AUCTION_ASSET_FROZEN AuctionState = "ASSET_FROZEN"
	AUCTION_PAYING       AuctionState = "PAYING"
	AUCTION_PAID         AuctionState = "PAID"
	AUCTION_NO_BIDS      AuctionState = "NO_BIDS"
)
