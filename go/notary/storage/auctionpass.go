package storage

type AuctionPass struct {
	Owner              string
	PurchasedForTeamId string
	ProductId          string
	Ack                bool
}

func NewAuctionPass() *AuctionPass {
	order := AuctionPass{}
	order.Ack = false
	return &order
}
