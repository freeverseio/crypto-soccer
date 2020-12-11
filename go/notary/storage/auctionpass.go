package storage

type AuctionPass struct {
	Owner string
}

func NewAuctionPass() *AuctionPass {
	order := AuctionPass{}
	return &order
}
