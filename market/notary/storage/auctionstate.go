package storage

type AuctionState string

const (
	STARTED      AuctionState = "STARTED"
	ASSET_FROZEN AuctionState = "ASSET_FROZEN"
	PAYING       AuctionState = "PAYING"
	PAID         AuctionState = "PAID"
	NO_BIDS      AuctionState = "NO_BIDS"
)
